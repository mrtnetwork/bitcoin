package example

import (
	"fmt"
	"github.com/mrtnetwork/bitcoin/address"
	"testing"
)

func ExampleAddress(t *testing.T) {

	// P2PKH ADDRESS
	// address in testnet:  myVMJgRi6arv4hLbeUcJYKUJWmFnpjtVme
	// address in mainnet:  1JyQ1dLjHZRfHaryvudviQFyemf5vjbmUf
	exampleAddr1, _ := address.P2PKHAddressFromAddress("myVMJgRi6arv4hLbeUcJYKUJWmFnpjtVme")
	fmt.Println("address in testnet: ", exampleAddr1.Show(address.TestnetNetwork))
	fmt.Println("address in mainnet: ", exampleAddr1.Show(address.MainnetNetwork))

	// P2TR ADDRESS
	// address in testnet:  tb1pyhmqwlcrws4dxcgalt4mrffgnys879vs59xf6sve4hazyvmhecxq3e6sc0
	// address in mainnet:  bc1pyhmqwlcrws4dxcgalt4mrffgnys879vs59xf6sve4hazyvmhecxqx3vlzq
	exampleAddr2, _ := address.P2TRAddressFromAddress("tb1pyhmqwlcrws4dxcgalt4mrffgnys879vs59xf6sve4hazyvmhecxq3e6sc0")
	fmt.Println("address in testnet: ", exampleAddr2.Show(address.TestnetNetwork))
	fmt.Println("address in mainnet: ", exampleAddr2.Show(address.MainnetNetwork))

	// P2SH(P2PKH) ADDRESS
	// address in testnet:  2N2yqygBJRvDzLzvPe91qKfSYnK5utGckJX
	// address in mainnet:  3BRduwFGpTie9DHqy1PxhiTHZxsk7jpr9x
	exampleAddr3, _ := address.P2SHAddressFromAddress("2N2yqygBJRvDzLzvPe91qKfSYnK5utGckJX", address.P2PKHInP2SH)
	fmt.Println("address in testnet: ", exampleAddr3.Show(address.TestnetNetwork))
	fmt.Println("address in mainnet: ", exampleAddr3.Show(address.MainnetNetwork))

	// P2PKH ADDRESS
	// address in testnet:  mzUzciYUGsNxLCaaHwou27F4RbnDTzKomV
	// address in mainnet:  1Ky3KfTVTqwhZ66xaNqXCC2jZcBWZD6ppQ
	exampleAddr4, _ := address.P2PKHAddressFromAddress("mzUzciYUGsNxLCaaHwou27F4RbnDTzKomV")
	fmt.Println("address in testnet: ", exampleAddr4.Show(address.TestnetNetwork))
	fmt.Println("address in mainnet: ", exampleAddr4.Show(address.MainnetNetwork))

	// P2SH(P2PKH) ADDRESS
	// address in testnet:  2MzibgEeJYCN8mjJsZTg79AH7au4PCkHXHo
	// address in mainnet:  39APcViGvjrnZwgKtL4EXDHrNYrDT9YuCQ
	exampleAddr5, _ := address.P2SHAddressFromAddress("2MzibgEeJYCN8mjJsZTg79AH7au4PCkHXHo", address.P2PKHInP2SH)
	fmt.Println("address in testnet: ", exampleAddr5.Show(address.TestnetNetwork))
	fmt.Println("address in mainnet: ", exampleAddr5.Show(address.MainnetNetwork))

	// P2SH(P2WSH) ADDRESS
	// address in testnet:  2N7bNV1WPwCVHfoqqRhvtmbAfktazfjHEW2
	// address in mainnet:  3G3ARGaNKjywU2DHkaK29eBQYYNppD2xDE
	exampleAddr6, _ := address.P2SHAddressFromAddress("2N7bNV1WPwCVHfoqqRhvtmbAfktazfjHEW2", address.P2WSHInP2SH)
	fmt.Println("address in testnet: ", exampleAddr6.Show(address.TestnetNetwork))
	fmt.Println("address in mainnet: ", exampleAddr6.Show(address.MainnetNetwork))

	// P2SH(P2WPKH) ADDRESS
	// address in testnet:  2N38S8G9q6qyEjPWmicqxrVwjq4QiTbcyf4
	// address in mainnet:  3BaE4XDoVPTtXbtE3VE6EYxUciCYd5rspb
	exampleAddr7, _ := address.P2SHAddressFromAddress("2N38S8G9q6qyEjPWmicqxrVwjq4QiTbcyf4", address.P2WPKHInP2SH)
	fmt.Println("address in testnet: ", exampleAddr7.Show(address.TestnetNetwork))
	fmt.Println("address in mainnet: ", exampleAddr7.Show(address.MainnetNetwork))

	// P2SH(P2PK) ADDRESS
	// address in testnet:  2MugsNcgzLJ1HosnZyC2CfZVmgbMPK1XubR
	// address in mainnet:  348fJskxiqVwc6A2J4QL3cWWUF9DbWjTni
	exampleAddr8, _ := address.P2SHAddressFromAddress("2MugsNcgzLJ1HosnZyC2CfZVmgbMPK1XubR", address.P2PKInP2SH)
	fmt.Println("address in testnet: ", exampleAddr8.Show(address.TestnetNetwork))
	fmt.Println("address in mainnet: ", exampleAddr8.Show(address.MainnetNetwork))

	// P2WPKH ADDRESS
	// address in testnet:  tb1q6q9halaazasd42gzsc2cvv5xls295w7kawkhxy
	// address in mainnet:  bc1q6q9halaazasd42gzsc2cvv5xls295w7khgdyah
	exampleAddr9, _ := address.P2WPKHAddresssFromAddress("tb1q6q9halaazasd42gzsc2cvv5xls295w7kawkhxy")
	fmt.Println("address in testnet: ", exampleAddr9.Show(address.TestnetNetwork))
	fmt.Println("address in mainnet: ", exampleAddr9.Show(address.MainnetNetwork))

}
