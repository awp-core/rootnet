// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package bindings

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// SubnetNFTSubnetData is an auto generated low-level Go binding around an user-defined struct.
type SubnetNFTSubnetData struct {
	Name          string
	SubnetManager common.Address
	AlphaToken    common.Address
	SkillsURI     string
	MinStake      *big.Int
	Owner         common.Address
}

// SubnetNFTMetaData contains all meta data concerning the SubnetNFT contract.
var SubnetNFTMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"name_\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"symbol_\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"rootNet_\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"approve\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"balanceOf\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"burn\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getAlphaToken\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getApproved\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getMinStake\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint128\",\"internalType\":\"uint128\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getSubnetData\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structSubnetNFT.SubnetData\",\"components\":[{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"subnetManager\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"alphaToken\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"skillsURI\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"minStake\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getSubnetManager\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isApprovedForAll\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"mint\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"name_\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"subnetManager_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"alphaToken_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"minStake_\",\"type\":\"uint128\",\"internalType\":\"uint128\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"name\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"ownerOf\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"rootNet\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"safeTransferFrom\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"safeTransferFrom\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setApprovalForAll\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"approved\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setBaseURI\",\"inputs\":[{\"name\":\"uri\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setMinStake\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"minStake_\",\"type\":\"uint128\",\"internalType\":\"uint128\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setSkillsURI\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"skillsURI_\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"supportsInterface\",\"inputs\":[{\"name\":\"interfaceId\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"symbol\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"tokenURI\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"transferFrom\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"Approval\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"approved\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ApprovalForAll\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"operator\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"approved\",\"type\":\"bool\",\"indexed\":false,\"internalType\":\"bool\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"MinStakeUpdated\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"minStake\",\"type\":\"uint128\",\"indexed\":false,\"internalType\":\"uint128\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"SkillsURIUpdated\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"skillsURI\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Transfer\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"ERC721IncorrectOwner\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC721InsufficientApproval\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ERC721InvalidApprover\",\"inputs\":[{\"name\":\"approver\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC721InvalidOperator\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC721InvalidOwner\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC721InvalidReceiver\",\"inputs\":[{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC721InvalidSender\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC721NonexistentToken\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"NotRootNet\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotTokenOwner\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"TokenNotExist\",\"inputs\":[]}]",
	Bin: "0x60a060405234620003265762001cd4803803806200001d816200032a565b9283398101606082820312620003265781516001600160401b03908181116200032657826200004e91850162000350565b906020928385015182811162000326576040916200006e91870162000350565b940151926001600160a01b0384168403620003265782519082821162000244575f54916001948584811c941680156200031b575b8385101462000225578190601f94858111620002c8575b50839085831160011462000264575f9262000258575b50505f19600383901b1c191690851b175f555b8551928311620002445783548481811c9116801562000239575b828210146200022557828111620001dd575b5080918311600114620001785750819293945f926200016c575b50505f19600383901b1c191690821b1790555b6080526040516119139081620003c18239608051818181610a8301528181610c1601528181610d800152610e2e0152f35b015190505f8062000128565b90601f19831695845f52825f20925f905b888210620001c55750508385969710620001ac575b505050811b0190556200013b565b01515f1960f88460031b161c191690555f80806200019e565b80878596829496860151815501950193019062000189565b845f52815f208380860160051c8201928487106200021b575b0160051c019085905b8281106200020f5750506200010e565b5f8155018590620001ff565b92508192620001f6565b634e487b7160e01b5f52602260045260245ffd5b90607f1690620000fc565b634e487b7160e01b5f52604160045260245ffd5b015190505f80620000cf565b90879350601f198316915f8052855f20925f5b87828210620002b1575050841162000298575b505050811b015f55620000e2565b01515f1960f88460031b161c191690555f80806200028a565b8385015186558b9790950194938401930162000277565b9091505f8052835f208580850160051c82019286861062000311575b918991869594930160051c01915b82811062000302575050620000b9565b5f8155859450899101620002f2565b92508192620002e4565b93607f1693620000a2565b5f80fd5b6040519190601f01601f191682016001600160401b038111838210176200024457604052565b919080601f84011215620003265782516001600160401b038111620002445760209062000386601f8201601f191683016200032a565b9281845282828701011162000326575f5b818110620003ac5750825f9394955001015290565b85810183015184820184015282016200039756fe6080604090808252600480361015610015575f80fd5b5f3560e01c91826301ffc9a71461126e5750816306fdde03146111c1578163081812fc14611189578163095ea7b3146110ae57816323b872dd1461109757816335c1b08114610daf578163405a0b0614610d6c57816342842e0e14610d2057816342966c6814610bf457816355f804b314610a3b5781636352211e14610a0c57816363a9bbe51461097057816370a082311461091c57816373f231e7146108e95781637c2f4cd61461073f578163854744ca1461060857816395d89b4114610534578163a22cb46514610499578163b88d4fde14610438578163c7bc8ec614610406578163c87b56dd146101a0578163e630cb961461016e575063e985e9c51461011d575f80fd5b3461016a578060031936011261016a5760209061013861131e565b610140611334565b9060018060a01b038091165f5260058452825f2091165f52825260ff815f20541690519015158152f35b5f80fd5b823461016a57602036600319011261016a57602091355f526007825260018060a01b036001825f200154169051908152f35b823461016a576020908160031936011261016a578235926101c0846118a3565b505f939081907a184f03e93ff9f4daa797ed6e38ed64bf6a1f01000000000000000082818110156103f7575b50506d04ee2d6d415b85acef8100000000808310156103e9575b50662386f26fc10000808310156103da575b506305f5e100808310156103cb575b50612710808310156103bf575b505060648110156103b1575b600a809110156103a7575b60018086019161027261025d846113ff565b9361026a875195866113dd565b8085526113ff565b9382602188860199601f19809801368c37860101905b610379575b50505083519586925f92600654906102a482611451565b91896001821691825f146103575750506001146102f7575b505091806102d76102e495936102f3989795519384916112d8565b01039081018652856113dd565b519282849384528301906112f9565b0390f35b889293945060065f527ff652222313e28459528d920b65115c16c04f3efc82aaedc97be59f3f377c0d3f905f915b83831061033d575050508401019190806102d76102bc565b80549783018501979097528a968a94909201918101610325565b60ff191681890152831515909302870190920194508291506102d790506102bc565b5f19019082906f181899199a1a9b1b9c1cb0b131b232b360811b8282061a835304908382610288575061028d565b936001019361024b565b606460029104940193610240565b95019490048580610234565b60089196920491019486610227565b60109196920491019486610218565b859196920491019486610206565b919650915004829486806101ec565b823461016a57602036600319011261016a57602091355f526007825260018060a01b036002825f200154169051908152f35b3461016a57608036600319011261016a5761045161131e565b610459611334565b9060643567ffffffffffffffff811161016a573660238201121561016a576104979381602461048d9336930135910161141b565b9160443591611757565b005b823461016a578060031936011261016a576104b261131e565b906024359182151580930361016a576001600160a01b031692831561051f5750335f526005602052805f20835f52602052805f2060ff1981541660ff8416179055519081527f17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c3160203392a3005b836024925191630b61174360e31b8352820152fd5b823461016a575f36600319011261016a578051905f90826001926001549361055b85611451565b90818452602095866001821691825f146105e657505060011461058b575b50506102f392916102e49103856113dd565b9085925060015f527fb10e2d527612073b26eecdfd717e6a320cf44b4afac2b0732d9fcbe2b7fa0cf6915f925b8284106105ce57505050820101816102e4610579565b8054848a0186015288955087949093019281016105b8565b60ff19168682015292151560051b850190920192508391506102e49050610579565b90503461016a576020918260031936011261016a5781355f60a0835161062d816113ad565b60608152828782015282858201526060808201528260808201520152805f526002845260018060a01b039081835f20541693841561073157505f5260078452815f209260088552825f2092826001860154168360028701541691846001600160801b039384600189015416978351996106a58b6113ad565b6106ae90611489565b8a528a8a01948552838a019182526106c590611489565b9260608a0193845260808a0198895260a08a019687528281519b8c9b818d5251908c0160c0905260e08c016106f9916112f9565b955116908a01525116606088015251868203601f1901608088015261071e91906112f9565b93511660a0850152511660c08301520390f35b835163224a1b1160e11b8152fd5b823461016a578060031936011261016a5781359067ffffffffffffffff9060243582811161016a57610774903690860161137f565b9261077e856118a3565b336001600160a01b03909116036108d957845f5260209560088752835f209185116108c657506107b8846107b28354611451565b836116d7565b5f601f851160011461083e577fd1332ed84c54e159e7a4245f8e021aff5b3389b685598c228394168ae30d1020959692859261080b848080955f91610833575b508160011b915f199060031b1c19161790565b90555b8451958487958652850152848401375f828201840152601f01601f19168101030190a2005b90508401358c6107f8565b601f19851690825f52875f20915f5b8181106108af5750928693927fd1332ed84c54e159e7a4245f8e021aff5b3389b685598c228394168ae30d10209899959285809510610896575b5050600183811b01905561080e565b8301355f19600386901b60f8161c191690558980610887565b91928960018192868901358155019401920161084d565b604190634e487b7160e01b5f525260245ffd5b82516359dc379f60e01b81528690fd5b823461016a57602036600319011261016a57602091355f52600882526001600160801b036001825f200154169051908152f35b823461016a57602036600319011261016a576001600160a01b0361093e61131e565b16801561095a57602092505f5260038252805f20549051908152f35b81516322718ad960e21b81525f81850152602490fd5b90503461016a578160031936011261016a57803591602435916001600160801b03831680930361016a576109a3846118a3565b336001600160a01b03909116036109fe57507fd0b53d029b5624436a948f7e5e2d9854defd8058cb4a20ff51a0ff9599ad6de891602091845f52600883526001815f2001826001600160801b031982541617905551908152a2005b90516359dc379f60e01b8152fd5b823461016a57602036600319011261016a57610a2a602092356118a3565b90516001600160a01b039091168152f35b90503461016a5760208060031936011261016a5767ffffffffffffffff91803583811161016a573660238201121561016a57610a80903690602481850135910161141b565b937f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03163303610be8575083519283116108c65750610ac7600654611451565b601f8111610b8f575b50602090601f8311600114610b1157508190610b01935f92610b06575b50508160011b915f199060031b1c19161790565b600655005b015190505f80610aed565b90601f1983169360065f527ff652222313e28459528d920b65115c16c04f3efc82aaedc97be59f3f377c0d3f925f905b868210610b775750508360019510610b5f575b505050811b01600655005b01515f1960f88460031b161c191690555f8080610b54565b80600185968294968601518155019501930190610b41565b610bd89060065f527ff652222313e28459528d920b65115c16c04f3efc82aaedc97be59f3f377c0d3f601f850160051c81019160208610610bde575b601f0160051c01906116c1565b5f610ad0565b9091508190610bcb565b516257aacb60e11b8152fd5b90503461016a576020908160031936011261016a578035906001600160a01b037f000000000000000000000000000000000000000000000000000000000000000081163303610d1257825f5260028452845f205416825f8215928315610ce1575b8282526002875287822080546001600160a01b03191690557fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef8280a4610ccc575f6001856008868680865260078252856002858220610cb38161170f565b828882015501558552528220610cc88161170f565b0155005b602491845191637e27328960e01b8352820152fd5b5f83815260046020526040902080546001600160a01b03191690558082526003875287822082198154019055610c55565b5083516257aacb60e11b8152fd5b823461016a57610d2f3661134a565b91835193602085019085821067ffffffffffffffff831117610d59576104979650525f8452611757565b604187634e487b7160e01b5f525260245ffd5b823461016a575f36600319011261016a57517f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03168152602090f35b90503461016a5760c036600319011261016a57610dca61131e565b916024356044359267ffffffffffffffff9384811161016a57610df0903690830161137f565b6001600160a01b0360643581811693929084900361016a576084359382851680950361016a5760a435996001600160801b038b16809b0361016a57837f000000000000000000000000000000000000000000000000000000000000000016330361108857831698891561107257875f5260209660028852848a5f2054169586151580611041575b8c5f5260038a528a8c5f209d60019e8f8154019055815f5260028c528d5f20996001600160601b0360a01b9a828c8254161790557fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef5f80a461102b578a51946060860186811085821117611018578c52610ef39136919061141b565b8452878401928352898401968752885f5260078852895f2093519081519283116108c65750610f2c82610f268654611451565b866116d7565b8790601f8311600114610fa7579180610f5e9260029695945f92610b065750508160011b915f199060031b1c19161790565b82555b838a830191511685825416179055019251169082541617905584610f8157005b5f91825260089052200180546fffffffffffffffffffffffffffffffff19169091179055005b9392918b91601f19821690855f528a5f20915f5b8c828210610ff5575050968360029810610fdd575b505050811b018255610f61565b01515f1960f88460031b161c191690555f8080610fd0565b919295899487849397999a9b015181550195019301908e94929796959391610fbb565b604184634e487b7160e01b5f525260245ffd5b8a516339e3563760e11b81525f81840152602490fd5b5f8b815260046020526040812080546001600160a01b031916905588815260038b528c902080545f19019055610e77565b8851633250574960e11b81525f81890152602490fd5b88516257aacb60e11b81528790fd5b3461016a576104976110a83661134a565b9161152b565b823461016a578060031936011261016a576110c761131e565b916024356110d4816118a3565b33151580611176575b8061114f575b611139576001600160a01b039485169482918691167f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b9255f80a45f526020525f20906001600160601b0360a01b8254161790555f80f35b835163a9fbf51f60e01b81523381850152602490fd5b5060018060a01b0381165f526005602052835f20335f5260205260ff845f205416156110e3565b506001600160a01b0381163314156110dd565b823461016a57602036600319011261016a5781602092356111a9816118a3565b505f52825260018060a01b03815f2054169051908152f35b823461016a575f36600319011261016a578051905f90825f54926111e484611451565b808352602094600190866001821691825f146105e65750506001146112155750506102f392916102e49103856113dd565b5f80805286935091907f290decd9548b62a8d60345a988386fc84ba6bc95484008f6362f93160ef3e5635b82841061125657505050820101816102e4610579565b8054848a018601528895508794909301928101611240565b903461016a57602036600319011261016a57359063ffffffff60e01b821680920361016a576020916380ac58cd60e01b81149081156112c7575b81156112b6575b5015158152f35b6301ffc9a760e01b149050836112af565b635b5e139f60e01b811491506112a8565b5f5b8381106112e95750505f910152565b81810151838201526020016112da565b90602091611312815180928185528580860191016112d8565b601f01601f1916010190565b600435906001600160a01b038216820361016a57565b602435906001600160a01b038216820361016a57565b606090600319011261016a576001600160a01b0390600435828116810361016a5791602435908116810361016a579060443590565b9181601f8401121561016a5782359167ffffffffffffffff831161016a576020838186019501011161016a57565b60c0810190811067ffffffffffffffff8211176113c957604052565b634e487b7160e01b5f52604160045260245ffd5b90601f8019910116810190811067ffffffffffffffff8211176113c957604052565b67ffffffffffffffff81116113c957601f01601f191660200190565b929192611427826113ff565b9161143560405193846113dd565b82948184528183011161016a578281602093845f960137010152565b90600182811c9216801561147f575b602083101461146b57565b634e487b7160e01b5f52602260045260245ffd5b91607f1691611460565b9060405191825f825461149b81611451565b908184526020946001916001811690815f1461150957506001146114cb575b5050506114c9925003836113dd565b565b5f90815285812095935091905b8183106114f15750506114c993508201015f80806114ba565b855488840185015294850194879450918301916114d8565b925050506114c994925060ff191682840152151560051b8201015f80806114ba565b6001600160a01b0391821692909183156116a957815f52602092600284528260409583875f2054169533151580611617575b50600290876115e6575b825f5260038152885f2060018154019055835f5252865f20816001600160601b0360a01b825416179055857fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef5f80a416928383036115c55750505050565b6064945051926364283d7b60e01b8452600484015260248301526044820152fd5b5f84815260046020526040812080546001600160a01b03191690558881526003825289902080545f19019055611567565b9192509080611668575b1561162f579084915f61155d565b86858761164c576024915190637e27328960e01b82526004820152fd5b604491519063177e802f60e01b82523360048301526024820152fd5b50338614801561168d575b806116215750845f52600481523384885f20541614611621565b50855f5260058152865f20335f52815260ff875f205416611673565b604051633250574960e11b81525f6004820152602490fd5b8181106116cc575050565b5f81556001016116c1565b9190601f81116116e657505050565b6114c9925f5260205f20906020601f840160051c83019310610bde57601f0160051c01906116c1565b6117198154611451565b9081611723575050565b81601f5f9311600114611734575055565b908083918252611753601f60208420940160051c8401600185016116c1565b5555565b919261176484838561152b565b813b611771575b50505050565b604051630a85bd0160e11b8082523360048301526001600160a01b039485166024830152604482019590955260806064820152602095939092169391908590829081906117c29060848301906112f9565b03815f885af15f9181611863575b5061182d575050503d5f14611826573d6117e9816113ff565b906117f760405192836113dd565b81523d5f8483013e5b8051928361182157604051633250574960e11b815260048101849052602490fd5b019050fd5b6060611800565b9193506001600160e01b03199091160361184b57505f80808061176b565b60249060405190633250574960e11b82526004820152fd5b9091508581813d831161189c575b61187b81836113dd565b8101031261016a57516001600160e01b03198116810361016a57905f6117d0565b503d611871565b5f818152600260205260409020546001600160a01b03169081156118c5575090565b60249060405190637e27328960e01b82526004820152fdfea26469706673582212206910e363fe95f03391e3c09ad8103b8f58410b6e3b88f1272f2d85bdadcc39d964736f6c63430008180033",
}

// SubnetNFTABI is the input ABI used to generate the binding from.
// Deprecated: Use SubnetNFTMetaData.ABI instead.
var SubnetNFTABI = SubnetNFTMetaData.ABI

// SubnetNFTBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use SubnetNFTMetaData.Bin instead.
var SubnetNFTBin = SubnetNFTMetaData.Bin

// DeploySubnetNFT deploys a new Ethereum contract, binding an instance of SubnetNFT to it.
func DeploySubnetNFT(auth *bind.TransactOpts, backend bind.ContractBackend, name_ string, symbol_ string, rootNet_ common.Address) (common.Address, *types.Transaction, *SubnetNFT, error) {
	parsed, err := SubnetNFTMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(SubnetNFTBin), backend, name_, symbol_, rootNet_)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &SubnetNFT{SubnetNFTCaller: SubnetNFTCaller{contract: contract}, SubnetNFTTransactor: SubnetNFTTransactor{contract: contract}, SubnetNFTFilterer: SubnetNFTFilterer{contract: contract}}, nil
}

// SubnetNFT is an auto generated Go binding around an Ethereum contract.
type SubnetNFT struct {
	SubnetNFTCaller     // Read-only binding to the contract
	SubnetNFTTransactor // Write-only binding to the contract
	SubnetNFTFilterer   // Log filterer for contract events
}

// SubnetNFTCaller is an auto generated read-only Go binding around an Ethereum contract.
type SubnetNFTCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SubnetNFTTransactor is an auto generated write-only Go binding around an Ethereum contract.
type SubnetNFTTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SubnetNFTFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type SubnetNFTFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SubnetNFTSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type SubnetNFTSession struct {
	Contract     *SubnetNFT        // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// SubnetNFTCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type SubnetNFTCallerSession struct {
	Contract *SubnetNFTCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts    // Call options to use throughout this session
}

// SubnetNFTTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type SubnetNFTTransactorSession struct {
	Contract     *SubnetNFTTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// SubnetNFTRaw is an auto generated low-level Go binding around an Ethereum contract.
type SubnetNFTRaw struct {
	Contract *SubnetNFT // Generic contract binding to access the raw methods on
}

// SubnetNFTCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type SubnetNFTCallerRaw struct {
	Contract *SubnetNFTCaller // Generic read-only contract binding to access the raw methods on
}

// SubnetNFTTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type SubnetNFTTransactorRaw struct {
	Contract *SubnetNFTTransactor // Generic write-only contract binding to access the raw methods on
}

// NewSubnetNFT creates a new instance of SubnetNFT, bound to a specific deployed contract.
func NewSubnetNFT(address common.Address, backend bind.ContractBackend) (*SubnetNFT, error) {
	contract, err := bindSubnetNFT(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &SubnetNFT{SubnetNFTCaller: SubnetNFTCaller{contract: contract}, SubnetNFTTransactor: SubnetNFTTransactor{contract: contract}, SubnetNFTFilterer: SubnetNFTFilterer{contract: contract}}, nil
}

// NewSubnetNFTCaller creates a new read-only instance of SubnetNFT, bound to a specific deployed contract.
func NewSubnetNFTCaller(address common.Address, caller bind.ContractCaller) (*SubnetNFTCaller, error) {
	contract, err := bindSubnetNFT(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SubnetNFTCaller{contract: contract}, nil
}

// NewSubnetNFTTransactor creates a new write-only instance of SubnetNFT, bound to a specific deployed contract.
func NewSubnetNFTTransactor(address common.Address, transactor bind.ContractTransactor) (*SubnetNFTTransactor, error) {
	contract, err := bindSubnetNFT(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SubnetNFTTransactor{contract: contract}, nil
}

// NewSubnetNFTFilterer creates a new log filterer instance of SubnetNFT, bound to a specific deployed contract.
func NewSubnetNFTFilterer(address common.Address, filterer bind.ContractFilterer) (*SubnetNFTFilterer, error) {
	contract, err := bindSubnetNFT(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SubnetNFTFilterer{contract: contract}, nil
}

// bindSubnetNFT binds a generic wrapper to an already deployed contract.
func bindSubnetNFT(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := SubnetNFTMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SubnetNFT *SubnetNFTRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SubnetNFT.Contract.SubnetNFTCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SubnetNFT *SubnetNFTRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SubnetNFT.Contract.SubnetNFTTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SubnetNFT *SubnetNFTRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SubnetNFT.Contract.SubnetNFTTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SubnetNFT *SubnetNFTCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SubnetNFT.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SubnetNFT *SubnetNFTTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SubnetNFT.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SubnetNFT *SubnetNFTTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SubnetNFT.Contract.contract.Transact(opts, method, params...)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_SubnetNFT *SubnetNFTCaller) BalanceOf(opts *bind.CallOpts, owner common.Address) (*big.Int, error) {
	var out []interface{}
	err := _SubnetNFT.contract.Call(opts, &out, "balanceOf", owner)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_SubnetNFT *SubnetNFTSession) BalanceOf(owner common.Address) (*big.Int, error) {
	return _SubnetNFT.Contract.BalanceOf(&_SubnetNFT.CallOpts, owner)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_SubnetNFT *SubnetNFTCallerSession) BalanceOf(owner common.Address) (*big.Int, error) {
	return _SubnetNFT.Contract.BalanceOf(&_SubnetNFT.CallOpts, owner)
}

// GetAlphaToken is a free data retrieval call binding the contract method 0xc7bc8ec6.
//
// Solidity: function getAlphaToken(uint256 tokenId) view returns(address)
func (_SubnetNFT *SubnetNFTCaller) GetAlphaToken(opts *bind.CallOpts, tokenId *big.Int) (common.Address, error) {
	var out []interface{}
	err := _SubnetNFT.contract.Call(opts, &out, "getAlphaToken", tokenId)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetAlphaToken is a free data retrieval call binding the contract method 0xc7bc8ec6.
//
// Solidity: function getAlphaToken(uint256 tokenId) view returns(address)
func (_SubnetNFT *SubnetNFTSession) GetAlphaToken(tokenId *big.Int) (common.Address, error) {
	return _SubnetNFT.Contract.GetAlphaToken(&_SubnetNFT.CallOpts, tokenId)
}

// GetAlphaToken is a free data retrieval call binding the contract method 0xc7bc8ec6.
//
// Solidity: function getAlphaToken(uint256 tokenId) view returns(address)
func (_SubnetNFT *SubnetNFTCallerSession) GetAlphaToken(tokenId *big.Int) (common.Address, error) {
	return _SubnetNFT.Contract.GetAlphaToken(&_SubnetNFT.CallOpts, tokenId)
}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(uint256 tokenId) view returns(address)
func (_SubnetNFT *SubnetNFTCaller) GetApproved(opts *bind.CallOpts, tokenId *big.Int) (common.Address, error) {
	var out []interface{}
	err := _SubnetNFT.contract.Call(opts, &out, "getApproved", tokenId)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(uint256 tokenId) view returns(address)
func (_SubnetNFT *SubnetNFTSession) GetApproved(tokenId *big.Int) (common.Address, error) {
	return _SubnetNFT.Contract.GetApproved(&_SubnetNFT.CallOpts, tokenId)
}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(uint256 tokenId) view returns(address)
func (_SubnetNFT *SubnetNFTCallerSession) GetApproved(tokenId *big.Int) (common.Address, error) {
	return _SubnetNFT.Contract.GetApproved(&_SubnetNFT.CallOpts, tokenId)
}

// GetMinStake is a free data retrieval call binding the contract method 0x73f231e7.
//
// Solidity: function getMinStake(uint256 tokenId) view returns(uint128)
func (_SubnetNFT *SubnetNFTCaller) GetMinStake(opts *bind.CallOpts, tokenId *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _SubnetNFT.contract.Call(opts, &out, "getMinStake", tokenId)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetMinStake is a free data retrieval call binding the contract method 0x73f231e7.
//
// Solidity: function getMinStake(uint256 tokenId) view returns(uint128)
func (_SubnetNFT *SubnetNFTSession) GetMinStake(tokenId *big.Int) (*big.Int, error) {
	return _SubnetNFT.Contract.GetMinStake(&_SubnetNFT.CallOpts, tokenId)
}

// GetMinStake is a free data retrieval call binding the contract method 0x73f231e7.
//
// Solidity: function getMinStake(uint256 tokenId) view returns(uint128)
func (_SubnetNFT *SubnetNFTCallerSession) GetMinStake(tokenId *big.Int) (*big.Int, error) {
	return _SubnetNFT.Contract.GetMinStake(&_SubnetNFT.CallOpts, tokenId)
}

// GetSubnetData is a free data retrieval call binding the contract method 0x854744ca.
//
// Solidity: function getSubnetData(uint256 tokenId) view returns((string,address,address,string,uint128,address))
func (_SubnetNFT *SubnetNFTCaller) GetSubnetData(opts *bind.CallOpts, tokenId *big.Int) (SubnetNFTSubnetData, error) {
	var out []interface{}
	err := _SubnetNFT.contract.Call(opts, &out, "getSubnetData", tokenId)

	if err != nil {
		return *new(SubnetNFTSubnetData), err
	}

	out0 := *abi.ConvertType(out[0], new(SubnetNFTSubnetData)).(*SubnetNFTSubnetData)

	return out0, err

}

// GetSubnetData is a free data retrieval call binding the contract method 0x854744ca.
//
// Solidity: function getSubnetData(uint256 tokenId) view returns((string,address,address,string,uint128,address))
func (_SubnetNFT *SubnetNFTSession) GetSubnetData(tokenId *big.Int) (SubnetNFTSubnetData, error) {
	return _SubnetNFT.Contract.GetSubnetData(&_SubnetNFT.CallOpts, tokenId)
}

// GetSubnetData is a free data retrieval call binding the contract method 0x854744ca.
//
// Solidity: function getSubnetData(uint256 tokenId) view returns((string,address,address,string,uint128,address))
func (_SubnetNFT *SubnetNFTCallerSession) GetSubnetData(tokenId *big.Int) (SubnetNFTSubnetData, error) {
	return _SubnetNFT.Contract.GetSubnetData(&_SubnetNFT.CallOpts, tokenId)
}

// GetSubnetManager is a free data retrieval call binding the contract method 0xe630cb96.
//
// Solidity: function getSubnetManager(uint256 tokenId) view returns(address)
func (_SubnetNFT *SubnetNFTCaller) GetSubnetManager(opts *bind.CallOpts, tokenId *big.Int) (common.Address, error) {
	var out []interface{}
	err := _SubnetNFT.contract.Call(opts, &out, "getSubnetManager", tokenId)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetSubnetManager is a free data retrieval call binding the contract method 0xe630cb96.
//
// Solidity: function getSubnetManager(uint256 tokenId) view returns(address)
func (_SubnetNFT *SubnetNFTSession) GetSubnetManager(tokenId *big.Int) (common.Address, error) {
	return _SubnetNFT.Contract.GetSubnetManager(&_SubnetNFT.CallOpts, tokenId)
}

// GetSubnetManager is a free data retrieval call binding the contract method 0xe630cb96.
//
// Solidity: function getSubnetManager(uint256 tokenId) view returns(address)
func (_SubnetNFT *SubnetNFTCallerSession) GetSubnetManager(tokenId *big.Int) (common.Address, error) {
	return _SubnetNFT.Contract.GetSubnetManager(&_SubnetNFT.CallOpts, tokenId)
}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
func (_SubnetNFT *SubnetNFTCaller) IsApprovedForAll(opts *bind.CallOpts, owner common.Address, operator common.Address) (bool, error) {
	var out []interface{}
	err := _SubnetNFT.contract.Call(opts, &out, "isApprovedForAll", owner, operator)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
func (_SubnetNFT *SubnetNFTSession) IsApprovedForAll(owner common.Address, operator common.Address) (bool, error) {
	return _SubnetNFT.Contract.IsApprovedForAll(&_SubnetNFT.CallOpts, owner, operator)
}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
func (_SubnetNFT *SubnetNFTCallerSession) IsApprovedForAll(owner common.Address, operator common.Address) (bool, error) {
	return _SubnetNFT.Contract.IsApprovedForAll(&_SubnetNFT.CallOpts, owner, operator)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_SubnetNFT *SubnetNFTCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _SubnetNFT.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_SubnetNFT *SubnetNFTSession) Name() (string, error) {
	return _SubnetNFT.Contract.Name(&_SubnetNFT.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_SubnetNFT *SubnetNFTCallerSession) Name() (string, error) {
	return _SubnetNFT.Contract.Name(&_SubnetNFT.CallOpts)
}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (_SubnetNFT *SubnetNFTCaller) OwnerOf(opts *bind.CallOpts, tokenId *big.Int) (common.Address, error) {
	var out []interface{}
	err := _SubnetNFT.contract.Call(opts, &out, "ownerOf", tokenId)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (_SubnetNFT *SubnetNFTSession) OwnerOf(tokenId *big.Int) (common.Address, error) {
	return _SubnetNFT.Contract.OwnerOf(&_SubnetNFT.CallOpts, tokenId)
}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (_SubnetNFT *SubnetNFTCallerSession) OwnerOf(tokenId *big.Int) (common.Address, error) {
	return _SubnetNFT.Contract.OwnerOf(&_SubnetNFT.CallOpts, tokenId)
}

// RootNet is a free data retrieval call binding the contract method 0x405a0b06.
//
// Solidity: function rootNet() view returns(address)
func (_SubnetNFT *SubnetNFTCaller) RootNet(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _SubnetNFT.contract.Call(opts, &out, "rootNet")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// RootNet is a free data retrieval call binding the contract method 0x405a0b06.
//
// Solidity: function rootNet() view returns(address)
func (_SubnetNFT *SubnetNFTSession) RootNet() (common.Address, error) {
	return _SubnetNFT.Contract.RootNet(&_SubnetNFT.CallOpts)
}

// RootNet is a free data retrieval call binding the contract method 0x405a0b06.
//
// Solidity: function rootNet() view returns(address)
func (_SubnetNFT *SubnetNFTCallerSession) RootNet() (common.Address, error) {
	return _SubnetNFT.Contract.RootNet(&_SubnetNFT.CallOpts)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_SubnetNFT *SubnetNFTCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _SubnetNFT.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_SubnetNFT *SubnetNFTSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _SubnetNFT.Contract.SupportsInterface(&_SubnetNFT.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_SubnetNFT *SubnetNFTCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _SubnetNFT.Contract.SupportsInterface(&_SubnetNFT.CallOpts, interfaceId)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_SubnetNFT *SubnetNFTCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _SubnetNFT.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_SubnetNFT *SubnetNFTSession) Symbol() (string, error) {
	return _SubnetNFT.Contract.Symbol(&_SubnetNFT.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_SubnetNFT *SubnetNFTCallerSession) Symbol() (string, error) {
	return _SubnetNFT.Contract.Symbol(&_SubnetNFT.CallOpts)
}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 tokenId) view returns(string)
func (_SubnetNFT *SubnetNFTCaller) TokenURI(opts *bind.CallOpts, tokenId *big.Int) (string, error) {
	var out []interface{}
	err := _SubnetNFT.contract.Call(opts, &out, "tokenURI", tokenId)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 tokenId) view returns(string)
func (_SubnetNFT *SubnetNFTSession) TokenURI(tokenId *big.Int) (string, error) {
	return _SubnetNFT.Contract.TokenURI(&_SubnetNFT.CallOpts, tokenId)
}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 tokenId) view returns(string)
func (_SubnetNFT *SubnetNFTCallerSession) TokenURI(tokenId *big.Int) (string, error) {
	return _SubnetNFT.Contract.TokenURI(&_SubnetNFT.CallOpts, tokenId)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address to, uint256 tokenId) returns()
func (_SubnetNFT *SubnetNFTTransactor) Approve(opts *bind.TransactOpts, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _SubnetNFT.contract.Transact(opts, "approve", to, tokenId)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address to, uint256 tokenId) returns()
func (_SubnetNFT *SubnetNFTSession) Approve(to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _SubnetNFT.Contract.Approve(&_SubnetNFT.TransactOpts, to, tokenId)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address to, uint256 tokenId) returns()
func (_SubnetNFT *SubnetNFTTransactorSession) Approve(to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _SubnetNFT.Contract.Approve(&_SubnetNFT.TransactOpts, to, tokenId)
}

// Burn is a paid mutator transaction binding the contract method 0x42966c68.
//
// Solidity: function burn(uint256 tokenId) returns()
func (_SubnetNFT *SubnetNFTTransactor) Burn(opts *bind.TransactOpts, tokenId *big.Int) (*types.Transaction, error) {
	return _SubnetNFT.contract.Transact(opts, "burn", tokenId)
}

// Burn is a paid mutator transaction binding the contract method 0x42966c68.
//
// Solidity: function burn(uint256 tokenId) returns()
func (_SubnetNFT *SubnetNFTSession) Burn(tokenId *big.Int) (*types.Transaction, error) {
	return _SubnetNFT.Contract.Burn(&_SubnetNFT.TransactOpts, tokenId)
}

// Burn is a paid mutator transaction binding the contract method 0x42966c68.
//
// Solidity: function burn(uint256 tokenId) returns()
func (_SubnetNFT *SubnetNFTTransactorSession) Burn(tokenId *big.Int) (*types.Transaction, error) {
	return _SubnetNFT.Contract.Burn(&_SubnetNFT.TransactOpts, tokenId)
}

// Mint is a paid mutator transaction binding the contract method 0x35c1b081.
//
// Solidity: function mint(address to, uint256 tokenId, string name_, address subnetManager_, address alphaToken_, uint128 minStake_) returns()
func (_SubnetNFT *SubnetNFTTransactor) Mint(opts *bind.TransactOpts, to common.Address, tokenId *big.Int, name_ string, subnetManager_ common.Address, alphaToken_ common.Address, minStake_ *big.Int) (*types.Transaction, error) {
	return _SubnetNFT.contract.Transact(opts, "mint", to, tokenId, name_, subnetManager_, alphaToken_, minStake_)
}

// Mint is a paid mutator transaction binding the contract method 0x35c1b081.
//
// Solidity: function mint(address to, uint256 tokenId, string name_, address subnetManager_, address alphaToken_, uint128 minStake_) returns()
func (_SubnetNFT *SubnetNFTSession) Mint(to common.Address, tokenId *big.Int, name_ string, subnetManager_ common.Address, alphaToken_ common.Address, minStake_ *big.Int) (*types.Transaction, error) {
	return _SubnetNFT.Contract.Mint(&_SubnetNFT.TransactOpts, to, tokenId, name_, subnetManager_, alphaToken_, minStake_)
}

// Mint is a paid mutator transaction binding the contract method 0x35c1b081.
//
// Solidity: function mint(address to, uint256 tokenId, string name_, address subnetManager_, address alphaToken_, uint128 minStake_) returns()
func (_SubnetNFT *SubnetNFTTransactorSession) Mint(to common.Address, tokenId *big.Int, name_ string, subnetManager_ common.Address, alphaToken_ common.Address, minStake_ *big.Int) (*types.Transaction, error) {
	return _SubnetNFT.Contract.Mint(&_SubnetNFT.TransactOpts, to, tokenId, name_, subnetManager_, alphaToken_, minStake_)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0x42842e0e.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId) returns()
func (_SubnetNFT *SubnetNFTTransactor) SafeTransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _SubnetNFT.contract.Transact(opts, "safeTransferFrom", from, to, tokenId)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0x42842e0e.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId) returns()
func (_SubnetNFT *SubnetNFTSession) SafeTransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _SubnetNFT.Contract.SafeTransferFrom(&_SubnetNFT.TransactOpts, from, to, tokenId)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0x42842e0e.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId) returns()
func (_SubnetNFT *SubnetNFTTransactorSession) SafeTransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _SubnetNFT.Contract.SafeTransferFrom(&_SubnetNFT.TransactOpts, from, to, tokenId)
}

// SafeTransferFrom0 is a paid mutator transaction binding the contract method 0xb88d4fde.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId, bytes data) returns()
func (_SubnetNFT *SubnetNFTTransactor) SafeTransferFrom0(opts *bind.TransactOpts, from common.Address, to common.Address, tokenId *big.Int, data []byte) (*types.Transaction, error) {
	return _SubnetNFT.contract.Transact(opts, "safeTransferFrom0", from, to, tokenId, data)
}

// SafeTransferFrom0 is a paid mutator transaction binding the contract method 0xb88d4fde.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId, bytes data) returns()
func (_SubnetNFT *SubnetNFTSession) SafeTransferFrom0(from common.Address, to common.Address, tokenId *big.Int, data []byte) (*types.Transaction, error) {
	return _SubnetNFT.Contract.SafeTransferFrom0(&_SubnetNFT.TransactOpts, from, to, tokenId, data)
}

// SafeTransferFrom0 is a paid mutator transaction binding the contract method 0xb88d4fde.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId, bytes data) returns()
func (_SubnetNFT *SubnetNFTTransactorSession) SafeTransferFrom0(from common.Address, to common.Address, tokenId *big.Int, data []byte) (*types.Transaction, error) {
	return _SubnetNFT.Contract.SafeTransferFrom0(&_SubnetNFT.TransactOpts, from, to, tokenId, data)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_SubnetNFT *SubnetNFTTransactor) SetApprovalForAll(opts *bind.TransactOpts, operator common.Address, approved bool) (*types.Transaction, error) {
	return _SubnetNFT.contract.Transact(opts, "setApprovalForAll", operator, approved)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_SubnetNFT *SubnetNFTSession) SetApprovalForAll(operator common.Address, approved bool) (*types.Transaction, error) {
	return _SubnetNFT.Contract.SetApprovalForAll(&_SubnetNFT.TransactOpts, operator, approved)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_SubnetNFT *SubnetNFTTransactorSession) SetApprovalForAll(operator common.Address, approved bool) (*types.Transaction, error) {
	return _SubnetNFT.Contract.SetApprovalForAll(&_SubnetNFT.TransactOpts, operator, approved)
}

// SetBaseURI is a paid mutator transaction binding the contract method 0x55f804b3.
//
// Solidity: function setBaseURI(string uri) returns()
func (_SubnetNFT *SubnetNFTTransactor) SetBaseURI(opts *bind.TransactOpts, uri string) (*types.Transaction, error) {
	return _SubnetNFT.contract.Transact(opts, "setBaseURI", uri)
}

// SetBaseURI is a paid mutator transaction binding the contract method 0x55f804b3.
//
// Solidity: function setBaseURI(string uri) returns()
func (_SubnetNFT *SubnetNFTSession) SetBaseURI(uri string) (*types.Transaction, error) {
	return _SubnetNFT.Contract.SetBaseURI(&_SubnetNFT.TransactOpts, uri)
}

// SetBaseURI is a paid mutator transaction binding the contract method 0x55f804b3.
//
// Solidity: function setBaseURI(string uri) returns()
func (_SubnetNFT *SubnetNFTTransactorSession) SetBaseURI(uri string) (*types.Transaction, error) {
	return _SubnetNFT.Contract.SetBaseURI(&_SubnetNFT.TransactOpts, uri)
}

// SetMinStake is a paid mutator transaction binding the contract method 0x63a9bbe5.
//
// Solidity: function setMinStake(uint256 tokenId, uint128 minStake_) returns()
func (_SubnetNFT *SubnetNFTTransactor) SetMinStake(opts *bind.TransactOpts, tokenId *big.Int, minStake_ *big.Int) (*types.Transaction, error) {
	return _SubnetNFT.contract.Transact(opts, "setMinStake", tokenId, minStake_)
}

// SetMinStake is a paid mutator transaction binding the contract method 0x63a9bbe5.
//
// Solidity: function setMinStake(uint256 tokenId, uint128 minStake_) returns()
func (_SubnetNFT *SubnetNFTSession) SetMinStake(tokenId *big.Int, minStake_ *big.Int) (*types.Transaction, error) {
	return _SubnetNFT.Contract.SetMinStake(&_SubnetNFT.TransactOpts, tokenId, minStake_)
}

// SetMinStake is a paid mutator transaction binding the contract method 0x63a9bbe5.
//
// Solidity: function setMinStake(uint256 tokenId, uint128 minStake_) returns()
func (_SubnetNFT *SubnetNFTTransactorSession) SetMinStake(tokenId *big.Int, minStake_ *big.Int) (*types.Transaction, error) {
	return _SubnetNFT.Contract.SetMinStake(&_SubnetNFT.TransactOpts, tokenId, minStake_)
}

// SetSkillsURI is a paid mutator transaction binding the contract method 0x7c2f4cd6.
//
// Solidity: function setSkillsURI(uint256 tokenId, string skillsURI_) returns()
func (_SubnetNFT *SubnetNFTTransactor) SetSkillsURI(opts *bind.TransactOpts, tokenId *big.Int, skillsURI_ string) (*types.Transaction, error) {
	return _SubnetNFT.contract.Transact(opts, "setSkillsURI", tokenId, skillsURI_)
}

// SetSkillsURI is a paid mutator transaction binding the contract method 0x7c2f4cd6.
//
// Solidity: function setSkillsURI(uint256 tokenId, string skillsURI_) returns()
func (_SubnetNFT *SubnetNFTSession) SetSkillsURI(tokenId *big.Int, skillsURI_ string) (*types.Transaction, error) {
	return _SubnetNFT.Contract.SetSkillsURI(&_SubnetNFT.TransactOpts, tokenId, skillsURI_)
}

// SetSkillsURI is a paid mutator transaction binding the contract method 0x7c2f4cd6.
//
// Solidity: function setSkillsURI(uint256 tokenId, string skillsURI_) returns()
func (_SubnetNFT *SubnetNFTTransactorSession) SetSkillsURI(tokenId *big.Int, skillsURI_ string) (*types.Transaction, error) {
	return _SubnetNFT.Contract.SetSkillsURI(&_SubnetNFT.TransactOpts, tokenId, skillsURI_)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 tokenId) returns()
func (_SubnetNFT *SubnetNFTTransactor) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _SubnetNFT.contract.Transact(opts, "transferFrom", from, to, tokenId)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 tokenId) returns()
func (_SubnetNFT *SubnetNFTSession) TransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _SubnetNFT.Contract.TransferFrom(&_SubnetNFT.TransactOpts, from, to, tokenId)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 tokenId) returns()
func (_SubnetNFT *SubnetNFTTransactorSession) TransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _SubnetNFT.Contract.TransferFrom(&_SubnetNFT.TransactOpts, from, to, tokenId)
}

// SubnetNFTApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the SubnetNFT contract.
type SubnetNFTApprovalIterator struct {
	Event *SubnetNFTApproval // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *SubnetNFTApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SubnetNFTApproval)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(SubnetNFTApproval)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *SubnetNFTApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SubnetNFTApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SubnetNFTApproval represents a Approval event raised by the SubnetNFT contract.
type SubnetNFTApproval struct {
	Owner    common.Address
	Approved common.Address
	TokenId  *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed approved, uint256 indexed tokenId)
func (_SubnetNFT *SubnetNFTFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, approved []common.Address, tokenId []*big.Int) (*SubnetNFTApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var approvedRule []interface{}
	for _, approvedItem := range approved {
		approvedRule = append(approvedRule, approvedItem)
	}
	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _SubnetNFT.contract.FilterLogs(opts, "Approval", ownerRule, approvedRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return &SubnetNFTApprovalIterator{contract: _SubnetNFT.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed approved, uint256 indexed tokenId)
func (_SubnetNFT *SubnetNFTFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *SubnetNFTApproval, owner []common.Address, approved []common.Address, tokenId []*big.Int) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var approvedRule []interface{}
	for _, approvedItem := range approved {
		approvedRule = append(approvedRule, approvedItem)
	}
	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _SubnetNFT.contract.WatchLogs(opts, "Approval", ownerRule, approvedRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SubnetNFTApproval)
				if err := _SubnetNFT.contract.UnpackLog(event, "Approval", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseApproval is a log parse operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed approved, uint256 indexed tokenId)
func (_SubnetNFT *SubnetNFTFilterer) ParseApproval(log types.Log) (*SubnetNFTApproval, error) {
	event := new(SubnetNFTApproval)
	if err := _SubnetNFT.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SubnetNFTApprovalForAllIterator is returned from FilterApprovalForAll and is used to iterate over the raw logs and unpacked data for ApprovalForAll events raised by the SubnetNFT contract.
type SubnetNFTApprovalForAllIterator struct {
	Event *SubnetNFTApprovalForAll // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *SubnetNFTApprovalForAllIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SubnetNFTApprovalForAll)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(SubnetNFTApprovalForAll)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *SubnetNFTApprovalForAllIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SubnetNFTApprovalForAllIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SubnetNFTApprovalForAll represents a ApprovalForAll event raised by the SubnetNFT contract.
type SubnetNFTApprovalForAll struct {
	Owner    common.Address
	Operator common.Address
	Approved bool
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterApprovalForAll is a free log retrieval operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: event ApprovalForAll(address indexed owner, address indexed operator, bool approved)
func (_SubnetNFT *SubnetNFTFilterer) FilterApprovalForAll(opts *bind.FilterOpts, owner []common.Address, operator []common.Address) (*SubnetNFTApprovalForAllIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _SubnetNFT.contract.FilterLogs(opts, "ApprovalForAll", ownerRule, operatorRule)
	if err != nil {
		return nil, err
	}
	return &SubnetNFTApprovalForAllIterator{contract: _SubnetNFT.contract, event: "ApprovalForAll", logs: logs, sub: sub}, nil
}

// WatchApprovalForAll is a free log subscription operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: event ApprovalForAll(address indexed owner, address indexed operator, bool approved)
func (_SubnetNFT *SubnetNFTFilterer) WatchApprovalForAll(opts *bind.WatchOpts, sink chan<- *SubnetNFTApprovalForAll, owner []common.Address, operator []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _SubnetNFT.contract.WatchLogs(opts, "ApprovalForAll", ownerRule, operatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SubnetNFTApprovalForAll)
				if err := _SubnetNFT.contract.UnpackLog(event, "ApprovalForAll", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseApprovalForAll is a log parse operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: event ApprovalForAll(address indexed owner, address indexed operator, bool approved)
func (_SubnetNFT *SubnetNFTFilterer) ParseApprovalForAll(log types.Log) (*SubnetNFTApprovalForAll, error) {
	event := new(SubnetNFTApprovalForAll)
	if err := _SubnetNFT.contract.UnpackLog(event, "ApprovalForAll", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SubnetNFTMinStakeUpdatedIterator is returned from FilterMinStakeUpdated and is used to iterate over the raw logs and unpacked data for MinStakeUpdated events raised by the SubnetNFT contract.
type SubnetNFTMinStakeUpdatedIterator struct {
	Event *SubnetNFTMinStakeUpdated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *SubnetNFTMinStakeUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SubnetNFTMinStakeUpdated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(SubnetNFTMinStakeUpdated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *SubnetNFTMinStakeUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SubnetNFTMinStakeUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SubnetNFTMinStakeUpdated represents a MinStakeUpdated event raised by the SubnetNFT contract.
type SubnetNFTMinStakeUpdated struct {
	TokenId  *big.Int
	MinStake *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterMinStakeUpdated is a free log retrieval operation binding the contract event 0xd0b53d029b5624436a948f7e5e2d9854defd8058cb4a20ff51a0ff9599ad6de8.
//
// Solidity: event MinStakeUpdated(uint256 indexed tokenId, uint128 minStake)
func (_SubnetNFT *SubnetNFTFilterer) FilterMinStakeUpdated(opts *bind.FilterOpts, tokenId []*big.Int) (*SubnetNFTMinStakeUpdatedIterator, error) {

	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _SubnetNFT.contract.FilterLogs(opts, "MinStakeUpdated", tokenIdRule)
	if err != nil {
		return nil, err
	}
	return &SubnetNFTMinStakeUpdatedIterator{contract: _SubnetNFT.contract, event: "MinStakeUpdated", logs: logs, sub: sub}, nil
}

// WatchMinStakeUpdated is a free log subscription operation binding the contract event 0xd0b53d029b5624436a948f7e5e2d9854defd8058cb4a20ff51a0ff9599ad6de8.
//
// Solidity: event MinStakeUpdated(uint256 indexed tokenId, uint128 minStake)
func (_SubnetNFT *SubnetNFTFilterer) WatchMinStakeUpdated(opts *bind.WatchOpts, sink chan<- *SubnetNFTMinStakeUpdated, tokenId []*big.Int) (event.Subscription, error) {

	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _SubnetNFT.contract.WatchLogs(opts, "MinStakeUpdated", tokenIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SubnetNFTMinStakeUpdated)
				if err := _SubnetNFT.contract.UnpackLog(event, "MinStakeUpdated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseMinStakeUpdated is a log parse operation binding the contract event 0xd0b53d029b5624436a948f7e5e2d9854defd8058cb4a20ff51a0ff9599ad6de8.
//
// Solidity: event MinStakeUpdated(uint256 indexed tokenId, uint128 minStake)
func (_SubnetNFT *SubnetNFTFilterer) ParseMinStakeUpdated(log types.Log) (*SubnetNFTMinStakeUpdated, error) {
	event := new(SubnetNFTMinStakeUpdated)
	if err := _SubnetNFT.contract.UnpackLog(event, "MinStakeUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SubnetNFTSkillsURIUpdatedIterator is returned from FilterSkillsURIUpdated and is used to iterate over the raw logs and unpacked data for SkillsURIUpdated events raised by the SubnetNFT contract.
type SubnetNFTSkillsURIUpdatedIterator struct {
	Event *SubnetNFTSkillsURIUpdated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *SubnetNFTSkillsURIUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SubnetNFTSkillsURIUpdated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(SubnetNFTSkillsURIUpdated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *SubnetNFTSkillsURIUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SubnetNFTSkillsURIUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SubnetNFTSkillsURIUpdated represents a SkillsURIUpdated event raised by the SubnetNFT contract.
type SubnetNFTSkillsURIUpdated struct {
	TokenId   *big.Int
	SkillsURI string
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterSkillsURIUpdated is a free log retrieval operation binding the contract event 0xd1332ed84c54e159e7a4245f8e021aff5b3389b685598c228394168ae30d1020.
//
// Solidity: event SkillsURIUpdated(uint256 indexed tokenId, string skillsURI)
func (_SubnetNFT *SubnetNFTFilterer) FilterSkillsURIUpdated(opts *bind.FilterOpts, tokenId []*big.Int) (*SubnetNFTSkillsURIUpdatedIterator, error) {

	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _SubnetNFT.contract.FilterLogs(opts, "SkillsURIUpdated", tokenIdRule)
	if err != nil {
		return nil, err
	}
	return &SubnetNFTSkillsURIUpdatedIterator{contract: _SubnetNFT.contract, event: "SkillsURIUpdated", logs: logs, sub: sub}, nil
}

// WatchSkillsURIUpdated is a free log subscription operation binding the contract event 0xd1332ed84c54e159e7a4245f8e021aff5b3389b685598c228394168ae30d1020.
//
// Solidity: event SkillsURIUpdated(uint256 indexed tokenId, string skillsURI)
func (_SubnetNFT *SubnetNFTFilterer) WatchSkillsURIUpdated(opts *bind.WatchOpts, sink chan<- *SubnetNFTSkillsURIUpdated, tokenId []*big.Int) (event.Subscription, error) {

	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _SubnetNFT.contract.WatchLogs(opts, "SkillsURIUpdated", tokenIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SubnetNFTSkillsURIUpdated)
				if err := _SubnetNFT.contract.UnpackLog(event, "SkillsURIUpdated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseSkillsURIUpdated is a log parse operation binding the contract event 0xd1332ed84c54e159e7a4245f8e021aff5b3389b685598c228394168ae30d1020.
//
// Solidity: event SkillsURIUpdated(uint256 indexed tokenId, string skillsURI)
func (_SubnetNFT *SubnetNFTFilterer) ParseSkillsURIUpdated(log types.Log) (*SubnetNFTSkillsURIUpdated, error) {
	event := new(SubnetNFTSkillsURIUpdated)
	if err := _SubnetNFT.contract.UnpackLog(event, "SkillsURIUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SubnetNFTTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the SubnetNFT contract.
type SubnetNFTTransferIterator struct {
	Event *SubnetNFTTransfer // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *SubnetNFTTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SubnetNFTTransfer)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(SubnetNFTTransfer)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *SubnetNFTTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SubnetNFTTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SubnetNFTTransfer represents a Transfer event raised by the SubnetNFT contract.
type SubnetNFTTransfer struct {
	From    common.Address
	To      common.Address
	TokenId *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 indexed tokenId)
func (_SubnetNFT *SubnetNFTFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address, tokenId []*big.Int) (*SubnetNFTTransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}
	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _SubnetNFT.contract.FilterLogs(opts, "Transfer", fromRule, toRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return &SubnetNFTTransferIterator{contract: _SubnetNFT.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 indexed tokenId)
func (_SubnetNFT *SubnetNFTFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *SubnetNFTTransfer, from []common.Address, to []common.Address, tokenId []*big.Int) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}
	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _SubnetNFT.contract.WatchLogs(opts, "Transfer", fromRule, toRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SubnetNFTTransfer)
				if err := _SubnetNFT.contract.UnpackLog(event, "Transfer", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseTransfer is a log parse operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 indexed tokenId)
func (_SubnetNFT *SubnetNFTFilterer) ParseTransfer(log types.Log) (*SubnetNFTTransfer, error) {
	event := new(SubnetNFTTransfer)
	if err := _SubnetNFT.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
