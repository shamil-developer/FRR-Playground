# Комментрий команды FRR по поводу внедрения gRPC

Разработчик команды FRR дает следующий комментрий по поводу того что не все покртыо gRPC

[Ссылка на оригинальный комментарий](https://github.com/FRRouting/frr/discussions/11323#discussioncomment-16531189)

> [!TIP]
> 
> ## Комментрий разработчика
> 
> ### Sharing PoC findings -- FRR 10.3 gRPC for BGP automation (what works, what doesn't)
> 
> Hi everyone,
> 
> I'm working on a network automation platform (NaaS) that needs programmatic BGP management via FRR. We ran a hands-on PoC with FRR 10.3 compiled from source with `--enable-grpc --enable-mgmtd` to evaluate what's feasible today. Sharing our findings in case it helps others in the same situation.
> 
> #### Environment
> 
> - FRR 10.3 compiled from source (Debian Bookworm)
> - mgmtd with `-M grpc:50051`
> - bgpd with `-M grpc:50052` (separate port)
> - Real BGP sessions (eBGP, VRF, BFD, prefix-lists, route-maps)
> 
> #### What works via gRPC (mgmtd:50051)
> 
> **GetCapabilities** returns 13 YANG modules: `frr-affinity-map`, `frr-backend`, `frr-filter`, `frr-interface`, `frr-route-map`, > `frr-routing`, `frr-staticd`, `frr-vrf`, `frr-zebra`, `frr-zebra-route-map`, and 3 IETF modules. No `frr-bgp`.
> 
> **Full CRUD for prefix-lists and route-maps** via `EditCandidate` + `Commit`:
> - Create: `CreateCandidate` -> `EditCandidate` (XPaths like `/frr-filter:lib/prefix-list[type='ipv4'][name='test']/entry[sequence='10']/action`) -> `Commit(ALL)` -- works perfectly
> - Read: `Get(CONFIG, "/frr-filter:lib")` returns structured JSON with all entries
> - Delete: `EditCandidate` with delete paths -> `Commit` -- works
> - Transaction IDs, rollback support, validation phase -- all functional
> 
> **Important caveat**: prefix-lists created via mgmtd gRPC do NOT appear in `vtysh show ip prefix-list` (which queries bgpd directly). mgmtd and bgpd maintain separate config stores.
> 
> #### What doesn't work
> 
> **bgpd with `-M grpc:50052`** loads the gRPC module but GetCapabilities returns only 5 modules: `frr-bgp-route-map`, `frr-filter`, `frr-interface`, `frr-route-map`, `frr-vrf`. Still no `frr-bgp`.
> 
> Any Get request with BGP-related XPaths (`/frr-routing:routing/control-plane-protocols/control-plane-protocol[type='frr-bgp:bgp']...`) returns `Data path not found`. mgmtd log confirms: `libyang: Unknown/non-implemented module "frr-bgp"`.
> 
> The YANG files exist in `/yang/` (frr-bgp.yang, frr-bgp-neighbor.yang, etc.) but the northbound callbacks are not wired in bgpd source code.
> 
> #### Our workaround
> 
> We ended up with a hybrid approach in our Go client:
> - **gRPC** for policy management (prefix-lists, route-maps) -- proper transactional CRUD
> - **`vtysh -c "show bgp ... json"`** for BGP state reads -- structured JSON output
> - **`vtysh -f config.conf`** for BGP config writes
> 
> This is essentially what SONiC does with frrcfgd (CONFIG_DB -> vtysh) and what Cumulus NVUE does (REST -> regenerate frr.conf).
> 
> #### Questions for maintainers
> 
> 1. Is there any timeline or active work on wiring `frr-bgp` northbound callbacks into mgmtd? Issue #5428 tracks the roadmap but bgpd has no assignee.
> 2. Would the community benefit from a Go library wrapping vtysh JSON + gRPC for the modules that work? We're considering open-sourcing ours.
> 3. Is the mgmtd/bgpd config store separation (prefix-lists visible in gRPC but not in vtysh) expected behavior or a bug?
> 
> Thanks for the great work on FRR. The gRPC infrastructure is solid -- the modules that are wired up work flawlessly. Looking forward to bgpd joining the northbound family someday.