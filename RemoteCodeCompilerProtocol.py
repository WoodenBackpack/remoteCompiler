import scapy.packet as pc
import scapy.fields as fs

class RCCP(pc.Packet):
    name = "RemoteCodeCompilerPacket"
    fields_desc = [
        fs.XByteField("CPP", 1)
    ]
p = RCCP(CPP=1)
