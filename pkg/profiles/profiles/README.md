## Free5gc + UERANSIM --- Multi-UPF configuration

Here we provide profiles which allow to configure free5gc+ueransim demo to work with multiple UPF's.

```bash
# Profiles applicable to both configurations
./common
```

### Scenario 1: UPF, SMF, UERANSIM

```bash
# Profiles which allow to configure 3 UPF: UPFB, UPF1 and UPF2.
# In this scenario we create two PDU Sessions - two uesimtunX interfaces are created
#     - All traffic sent via uesimtun0 will go through UPFB -> UPF1 -> DN1
#     - All traffic sent via uesimtun1 will go through UPFB -> UPF2 -> DN2
./multi-pdu-session
```

### Scenario 2: UPF, SMF, UERANSIM

```bash
# Profiles which allow to configure 3 UPF: UPFB, UPF1 and UPF2
# In this scenario we create just one PDU Session - only uesimtun0 interface is created
# Uplink Classifier (ULCL) is used to forward traffic to the specific location
#     - By default all traffic will go through UPFB -> UPF2 -> DN
#     - If destination address is set to 192.168.101.6/32 traffic will go thorugh -> UPFB -> UPF1 -> DN
#     - If destination address is set to 192.168.101.4/32 traffic will go through -> UPFB -> UPF2 -> DN
# Specific maths are defined in smf-profiles as `ueRoutingInfo`
./ulcl-single-pdu-session
```

## Notes

* Please NOTE that ueRoutingInfo adjustment (for now) requires to completely restart free5gc. It's 
not possible to configure specificPaths on the fly [(implementation constraints).](https://forum.free5gc.org/t/while-running-in-ulcl-configuration-is-it-possible-to-change-the-uerouting-behaviour-without-restarting-the-free5gc-core/1051)
* For better understanding of ULCL please see [section 5.6.4.2 of ETSI TI 123 501 v15.4.0](https://www.etsi.org/deliver/etsi_ts/123500_123599/123501/15.04.00_60/ts_123501v150400p.pdf)
* Default UPF can't be set in the configuration files, it's hardcoded in the code - [check free5gc forum threat](https://forum.free5gc.org/t/how-to-specify-which-upf-is-used-as-a-default/969)
* This configuration applies to [towards5gs helm charts](https://github.com/Orange-OpenSource/towards5gs-helm/) - UERANSIM v3.2.4 & Free5gc v3.0.6

## How to package profiles

To create `*.tar.gz file` you can follow example below.

If you want to create `upf-profile.tar.gz` then go to `<demo-path>/pkg/profiles/profiles/ulcl-single-pdu-session/upf-profile/` and type

```bash
tar -czvf upf-profile.tar.gz override_values.yaml manifest.yaml
```

Then you can copy created file to the `<demo-path>/pkg/profiles`

```bash
cp upf-profile.tar.gz ../../../
```
