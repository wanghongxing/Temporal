package relay

import (
	"context"
	"fmt"
	"math/rand"

	pstore "gx/ipfs/QmQAGG1zxfePqj2t7bLxyN8AFccZ889DDR9Gn8kVLDrGZo/go-libp2p-peerstore"
	ma "gx/ipfs/QmRKLtwMw131aK7ugC3G7ybpumMz78YrJe5dzneyindvG1/go-multiaddr"
	peer "gx/ipfs/QmcqU6QUDSXprb1518vYDGczrTJTyGwLG9eUa5iNX4xUtS/go-libp2p-peer"
	tpt "gx/ipfs/QmdQx4ZhKGdv9TvpCFpMxFzjTQFHRmFqjBxkRVwzT1JNes/go-libp2p-transport"
)

func (d *RelayTransport) Dial(ctx context.Context, a ma.Multiaddr, p peer.ID) (tpt.Conn, error) {
	c, err := d.Relay().Dial(ctx, a, p)
	if err != nil {
		return nil, err
	}
	return d.upgrader.UpgradeOutbound(ctx, d, c, p)
}

func (r *Relay) Dial(ctx context.Context, a ma.Multiaddr, p peer.ID) (*Conn, error) {
	if !r.Matches(a) {
		return nil, fmt.Errorf("%s is not a relay address", a)
	}
	parts := ma.Split(a)

	spl := ma.Cast(ma.CodeToVarint(P_CIRCUIT))

	var relayaddr, destaddr ma.Multiaddr
	for i, p := range parts {
		if p.Equal(spl) {
			relayaddr = ma.Join(parts[:i]...)
			destaddr = ma.Join(parts[i+1:]...)
			break
		}
	}

	dinfo := &pstore.PeerInfo{ID: p, Addrs: []ma.Multiaddr{}}
	if len(destaddr.Bytes()) > 0 {
		dinfo.Addrs = append(dinfo.Addrs, destaddr)
	}

	if len(relayaddr.Bytes()) == 0 {
		// unspecific relay address, try dialing using known hop relays
		return r.tryDialRelays(ctx, *dinfo)
	}

	var rinfo *pstore.PeerInfo
	rinfo, err := pstore.InfoFromP2pAddr(relayaddr)
	if err != nil {
		return nil, fmt.Errorf("error parsing multiaddr '%s': %s", relayaddr.String(), err)
	}

	return r.DialPeer(ctx, *rinfo, *dinfo)
}

func (r *Relay) tryDialRelays(ctx context.Context, dinfo pstore.PeerInfo) (*Conn, error) {
	var relays []peer.ID
	r.mx.Lock()
	for p := range r.relays {
		relays = append(relays, p)
	}
	r.mx.Unlock()

	// shuffle list of relays, avoid overloading a specific relay
	for i := range relays {
		j := rand.Intn(i + 1)
		relays[i], relays[j] = relays[j], relays[i]
	}

	for _, relay := range relays {
		if len(r.host.Network().ConnsToPeer(relay)) == 0 {
			continue
		}

		rctx, cancel := context.WithTimeout(ctx, HopConnectTimeout)
		c, err := r.DialPeer(rctx, pstore.PeerInfo{ID: relay}, dinfo)
		cancel()

		if err == nil {
			return c, nil
		}

		log.Debugf("error opening relay connection through %s: %s", dinfo.ID, err.Error())
	}

	return nil, fmt.Errorf("Failed to dial through %d known relay hosts", len(relays))
}
