// Implementation of Bitcoin script handling, transaction creation, and Bitcoin script features.
package scripts

import (
	"github.com/mrtnetwork/bitcoin/constant"
	"github.com/mrtnetwork/bitcoin/digest"
	"github.com/mrtnetwork/bitcoin/formating"

	"encoding/binary"
	"encoding/hex"
	"fmt"
)

type ScriptType int

const (
	P2PKH ScriptType = iota
	P2SH
	P2WPKH
	P2WSH
	P2PK
)

// A Script contains just a list of OP_CODES and also knows how to serialize
// into bytes
type Script struct {
	// the list with all the script OP_CODES and data
	Script []interface{}
}

func NewScript(args ...interface{}) *Script {
	return &Script{Script: args}
}

func NewScriptFromList(args []interface{}) *Script {
	return &Script{Script: args}
}

func (s *Script) ToTapleafTaggedHash() []byte {
	leafVarBytes := []byte{constant.LEAF_VERSION_TAPSCRIPT}
	leafVarBytes = append(leafVarBytes, formating.PrependVarint(s.ToBytes())...)
	return digest.TaggedHash(leafVarBytes, "TapLeaf")
}

// Converts script to p2sh scriptPubKey (locking script)
// Calculates the hash160 (via the address) of the script and uses it to
// construct a P2SH script.
func (s *Script) ToP2shScriptPubKey() *Script {
	toBytes := s.ToBytes()
	h160 := digest.Hash160(toBytes)
	return NewScript("OP_HASH160", formating.BytesToHex(h160), "OP_EQUAL")
}

// Imports a Script commands list from raw hexadecimal data
func ScriptFromRaw(hexData string, hasSegwit bool) (*Script, error) {
	var commands []interface{}
	index := 0
	scriptraw, err := formating.HexToBytesCatch(hexData)
	if err != nil {
		return nil, fmt.Errorf("invalid script bytes")
	}

	for index < len(scriptraw) {
		b := int(scriptraw[index])
		if constant.CODE_OPS[b] != "" {
			commands = append(commands, constant.CODE_OPS[b])
			index++
		} else if !hasSegwit && b == 0x4c {
			bytesToRead := int(scriptraw[index+1])
			index++
			data := scriptraw[index : index+bytesToRead]
			commands = append(commands, hex.EncodeToString(data))
			index += bytesToRead
		} else if !hasSegwit && b == 0x4d {
			bytesToRead := int(binary.LittleEndian.Uint16(scriptraw[index+1 : index+3]))
			index += 3
			data := scriptraw[index : index+bytesToRead]
			commands = append(commands, hex.EncodeToString(data))
			index += bytesToRead
		} else if !hasSegwit && b == 0x4e {
			bytesToRead := int(binary.LittleEndian.Uint32(scriptraw[index+1 : index+5]))
			index += 5
			data := scriptraw[index : index+bytesToRead]
			commands = append(commands, hex.EncodeToString(data))
			index += bytesToRead
		} else {
			vi, size := formating.ViToInt(scriptraw[index:])
			dataSize := vi
			// size := size
			lastIndex := index + size + dataSize
			if lastIndex > len(scriptraw) {
				lastIndex = len(scriptraw)
			}
			commands = append(commands, hex.EncodeToString(scriptraw[index+size:lastIndex]))
			index += dataSize + size
		}
	}
	return NewScript(commands...), nil
}

// GetScriptType determines the script type based on the provided hash and whether it has
// SegWit data. It returns the identified ScriptType.
func GetScriptType(hash string, hasSegwit bool) (ScriptType, error) {
	s, err := ScriptFromRaw(hash, hasSegwit)
	if err != nil {
		return -1, err
	}
	if len(s.Script) == 0 {
		return -1, fmt.Errorf("invalid script bytes")
	}
	first := s.Script[0]
	sec := ""
	if len(s.Script) > 1 {
		sec = fmt.Sprintf("%v", s.Script[1])
	}
	th := ""
	if len(s.Script) > 2 {
		th = fmt.Sprintf("%v", s.Script[2])
	}
	four := ""
	if len(s.Script) > 3 {
		four = fmt.Sprintf("%v", s.Script[3])
	}
	five := ""
	if len(s.Script) > 4 {
		five = fmt.Sprintf("%v", s.Script[4])
	}

	if first == "OP_0" {
		if len(sec) == 40 {
			return P2WPKH, nil
		} else if len(sec) == 64 {
			return P2WSH, nil
		}
	} else if first == "OP_DUP" {
		if sec == "OP_HASH160" && four == "OP_EQUALVERIFY" && five == "OP_CHECKSIG" {
			return P2PKH, nil
		}
	} else if first == "OP_HASH160" && th == "OP_EQUAL" {
		return P2SH, nil
	} else if sec == "OP_CHECKSIG" && first != nil {
		if len(fmt.Sprintf("%v", first)) == 66 {
			return P2PK, nil
		}
	}
	return -1, fmt.Errorf("invalid script bytes")
}

// Converts the script to bytes
// If an OP code the appropriate byte is included according to:
// https://en.bitcoin.it/wiki/Script
// If not consider it data (signature, public key, public key hash, etc.) and
// and include with appropriate OP_PUSHDATA OP code plus length
func (s *Script) ToBytes() []byte {
	var scriptBytes []byte
	for _, token := range s.Script {
		if tokenStr, ok := token.(string); ok {
			if opCode, ok := constant.OP_CODES[tokenStr]; ok {
				scriptBytes = append(scriptBytes, opCode...)
				continue
			}
		}
		if tokenInt, ok := token.(int); ok && tokenInt >= 0 && tokenInt <= 16 {
			opCode := constant.OP_CODES[fmt.Sprintf("OP_%d", tokenInt)]
			scriptBytes = append(scriptBytes, opCode...)
		} else {
			if tokenInt, ok := token.(int); ok {
				scriptBytes = append(scriptBytes, formating.PushInteger(tokenInt)...)
			} else if tokenStr, ok := token.(string); ok {
				scriptBytes = append(scriptBytes, formating.OpPushData(tokenStr)...)
			}
		}
	}
	return scriptBytes
}

// returns a serialized version of the script in hex
func (s *Script) ToHex() string {
	bytes := s.ToBytes()
	return hex.EncodeToString(bytes)
}

// returns the list of strings that makes up this script
func ToScript(scrips []string) *Script {
	var interfaceList []interface{}
	for _, str := range scrips {
		interfaceList = append(interfaceList, str)
	}
	return &Script{Script: interfaceList}
}
