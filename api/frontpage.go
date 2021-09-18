package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"codeberg.org/imabritishcow/librarian/types"
	"codeberg.org/imabritishcow/librarian/utils"
	"github.com/dustin/go-humanize"
	"github.com/patrickmn/go-cache"
	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
)

var fpCache = cache.New(30*time.Minute, 30*time.Minute)

func GetFrontpageVideos() []types.Video {
	cacheData, found := fpCache.Get("fp")
	if found {
		return cacheData.([]types.Video)
	}

	claimSearchData := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "claim_search",
		"params": map[string]interface{}{
			"page_size":  20,
			"no_totals":  true,
			"claim_type": "stream",
			"any_tags":   []string{},
			"not_tags":   []string{"porn", "porno", "nsfw", "mature", "xxx", "sex", "creampie", "blowjob", "handjob", "vagina", "boobs", "big boobs", "big dick", "pussy", "cumshot", "anal", "hard fucking", "ass", "fuck", "hentai"},
			"channel_ids": []string{
				"c93267be27d3670fe09f9212e9412c9da432caec", "44c49bbab8a3e9999f3ba9dee0186288b1d960a7", "d78ea96e797cc5dbdc3162f1829d779f1e4d0d2a", "37f3125b8ae2a19082300e05951d6c33d65b84d3", "638f9e7b024efc24a14c161f8ccf4bacbcfda10b", "6184648aab0431c4c95c649072d1f9ff08b9bb7c", "b5d31cde873073718c033076656a27471e392afc",
				"b516294f541a18ce00b71a60b2c82ad2f87ff78d", "91e42cc450075f2c4c245bac7617bf903f16b4ce", "b6e207c5f8c58e7c8362cd05a1501bf2f5b694f2", "25f384bd95e218f6ac37fcaca99ed40f36760d8c", "f33657a2fcbab2dc3ce555d5d6728f8758af7bc7", "294f5c164da5ac9735658b2d58d8fee6745dfc45", "119a2e8c0b50f78d3861636d37c3b44ba8e689b5",
				"7b23cca3f49059f005e812be03931c81272eaac4", "fb0efeaa3788d1292bb49a94d77622503fe08129", "bc490776f367b8afccf0ea7349d657431ba1ded6", "48c7ea8bc2c4adba09bf21a29689e3b8c2967522", "df961194a798cc76306b9290701130c592530fb6", "cf0be9078d76951e2e228df68b5b0bbf71313aaa", "d746ac8d782f94d12d176c7a591f5bf8365bef3d",
				"4ad942982e43326c7700b1b6443049b3cfd82161", "1cdb5d0bdcb484907d0a2fea4efdfe0153838642", "6616707e1109aaa1c11b9f399f914d0cfb4f5303", "b7d02b4a0036114732c072269adb891dc5e34ca4", "9c51c1a119137cd17ed5ae09daa80c1cab6ac01d", "5f2a5c14b971a6f5eed0a67dc7af3a3fe5c0b6a4", "0e2b5b4cf59e859860000ff123dc12a317ad416b",
				"3fe68ad3da93065e35c37b14fbeef88b4b7785ed", "fd7ffcbafb74412a8812df4720feaf11fe70fe12", "92c0f2f3239f1f61496997bd2cdc197ec51bd423", "29193e9240a71a735639c66ee954e68414f11236", "25f384bd95e218f6ac37fcaca99ed40f36760d8c", "87b13b074936b1f42a7c6758c7c2995f58c602e7", "8d935c6c30510e1dfc10f803a9646fa8aa128b07",
				"25f384bd95e218f6ac37fcaca99ed40f36760d8c", "87b13b074936b1f42a7c6758c7c2995f58c602e7", "8d935c6c30510e1dfc10f803a9646fa8aa128b07", "8f4fecfc836ea33798ee3e5cef56926fa54e2cf9", "9a5dfcb1a4b29c3a1598392d039744b9938b5a26", "b39833be3032bbe1005f4f719f379a4621faeb13", "fb364ef587872515f545a5b4b3182b58073f230f",
				"b924ac36b7499591f7929d9c4903de79b07b1cb9", "72f9815b087b6d346745e3de71a6ce5fe73a8677", "b0198a465290f065378f3535666bee0653d6a9bb", "020ebeb40642bfb4bc3d9f6d28c098afc0a47481", "5c15e604c4207f52c8cf58fe21e63164c230e257", "273a2fa759f1a9f56b078633ea2f08fc2406002a", "930fc43ca7bae20d4706543e97175d1872b0671f",
				"0cb2ec46f06ba85520a1c1a56706acf35d5176dd", "057053dfb657aaa98553e2c544b06e1a2371557e", "64e091964a611a48424d254a3de2b952d0d6565a", "50ebba2b06908f93d7963b1c6826cc0fd6104477", "374ff82251a384601da73f30485c3ac8d7f4176b", "1487afc813124abbeb0629d2172be0f01ccec3bf", "b8bcc8bd8a9e7b3ba4ed074632601b377f9be387",
				"6a4fa1a68b92336e64006a4310cb160b07854329", "15f986a262fc6eff5774050c94d174c0533d505d", "1efa9b640ad980b2ec53834d60e9cff9554979cd", "064d4999ea15e433e06f16f391922390acab01cb", "4884e30b93b3c4c123a83154516196095f9e831e", "2827bfc459c12d7c6d280cbacee750811291d4ba", "273163260bceb95fa98d97d33d377c55395e329a",
				"9626816275585ac3443e7cddd1272c8652c23f1d", "a2e1bb1fed32c6a6290b679785dd31ca5c59cb5f", "d9535951222dd7a1ff7f763872cb0df44f7962bf", "243b6f18093ff97c861d0568c7d3379606201a4b", "1ce5ac7bab7f2e82af02305ced3c095e320b52e5", "3e63119a8503a6f666b0a736c8fdeb9e79d11eb4", "1cdb5d0bdcb484907d0a2fea4efdfe0153838642",
				"e33372c0d8b2cdd3e12252962ee1671d66143075", "7364ba2ac9090e468855ce9074bb50c306196f4c", "d6350f9158825662b99e4b5e0442bcc94d39bc11", "2a294ea41312c6da8de5ebd0f18dbfb2963bb1a4", "2305db36455aa0d18571015b9e9bd0950262aa0f", "82f1d8c257d3e76b711a5cecd1e49bd3fa6a9de9", "39ac239b5687f7d1c2ba74cd020b3547545dfdaf",
				"5737eb22f6119f27c9afccfe73ba710afd885371", "3583f5a570af6870504eea5a5f7afad6e1508508", "b0e489f986c345aef23c4a48d91cbcf5a6fdb9ac", "ba79c80788a9e1751e49ad401f5692d86f73a2db", "c54323436f50c633c870298bb374ac8e7560e6cd", "83725c7ee23bd4a8ca28a4fab0e313409def1dc7", "61c4e8636704f2f38bbe88b1f30ef0d74d6c0f49",
				"87ef9ba36019f7f3bf217cf47511645893b13f2e", "1bd0adfcf1c75bafc1ba3fc9b65a1a0470df6a91", "1be15348c51955179b7bf9aa90230a9425927ef6", "1f45ab3495df2c3be7d9c882bca9966305115cbb", "3346a4ff80b70ee7eea8a79fc79f19c43bb4464a", "452916a904030d2ce4ea2440fad2d0774e7296d9", "2675ef3adf52ebf8a44ff5da4306c293dfa6f901",
				"e2a76643735f8611a561794509c6bb2aac70eb04", "89c4cf244099918b1d3ed413df27d4216e97b499", "4b2b5822c8af3074c6ef9b789a8142d0ef623402", "3c5794f775975669745c412b0c30f48991d9e455", "1df464dbb302ced815c61431a5548a273e6de8e1", "1c0129135a23ff8096d533e0241921859f065f35", "6d2dac42692c279aed436f000c053c34c3ccec29",
				"403cd55421118093e83d490cf7991ebb477a320c", "b48a89b235fa577fab38ef56b57008c2aa0655fd", "06f6e7b9a05f12823c7d93dc04826636ae99c938", "e5872eb7237883e4158cf88e96a465f7a674c968", "b4ccefa53e0bb5730b56facb9b2ae5ad92c99935", "a30c5c8cb8faf78b7ad82118e8fd7039901a5a32", "5701b00a32d91035baa53e6a65800f0f5d4e520c",
				"6f061a2d8cf613a1ccca86468cf38143d2887c18", "4de30a6590a288684ab200409988c47c17a12ef2", "d0dec51ee295fae711466a6642749217b6e48032", "c73409178a35717a9a735aea942945fca53604fc", "f58e50a4c5fd2db81e1b5f4f5c36642d45f805c4", "074973ec6f89e0445903c844d02f9b15a64b2222", "f74655b4911e09c77d2fa0e700d5a76d2bc68ba5",
				"16fb96ca5b1713a4cc787e9e2542da724a45c854", "c1296488bc1aa7a8c8890509fd231267b3cf9e16", "0b7925083d4c382cb441f5f44ad626f7627bc198", "af0f5d9e055d07c8308da7f6726562420ff0bb82", "d6d17103b4aadd46c47236431eb7a96427ef98d7", "638f9e7b024efc24a14c161f8ccf4bacbcfda10b", "44492321a13a0087c12394d4fb48c208ca6d4b27",
				"cc8f04f02b15a4de03381f1ef506bd270c7b0936", "1698fe7e7cc0e748d99408e67df73f584e18113f", "38a83de67a246c72cd793104073968d7440e3509", "4c4c9e44dd32e618ce56c7ca2e2e2b53e8027779", "617886394bee35c359766fb4d9ae889b8fd39079", "88f9585ed0139d8775c47d2502b60ff028acc17e", "37f3125b8ae2a19082300e05951d6c33d65b84d3",
				"1dd1f89c9fec0450b8f446d88c8ca7eae806c54e", "7f745cf98bead9298afa74dc16fe736e69c74ba0", "1226c8d64dc52014a7fffde98d28acaab1969ce3", "76283e4cd0168e392128124a013c8335a987186d", "3f85de98c22a9ea0407c6421ba43f21701768bd9", "26cb3bd2123c671ebc24ccb37c87ed93047e7dfb", "d0754f60f7de56470372f876ee52d13a44019c85",
				"3151569ef709a0ec2721116a21028874e004a641", "3eb0394271d1c052cd0af482a304c29c00562bde", "1bd8bf2d2037a1bdd9bfdb1c4f767996389bd0d3", "5a638f268c78a312fe69fc1899bf8735bc6f5b73", "5b00937848d7c5f9287c38bab30a1c5c87f350a2", "e464c78688fdc5ac60dc5100fbf68f2fcc32b390", "2f2ea51d9117ab339dc72d5e25d051f873bf03fc",
				"ce605b7901923714aa0ea082ada13b0e8b099df6", "d7cfed7c3c9944c714bd79429a8d0ffff8056640", "ce605b7901923714aa0ea082ada13b0e8b099df6", "95d0f01df499867495cacac971eb9fdd0b4bc2e4", "9d40b87a50b576183d95ee2642faa4993a2748d6", "b797dc8d288818347fed3888b52e1abfb1368496", "277f7d20007960bee3e02774c61e7fb7481e0b14",
				"7b1301d743f9bf4db4445db1d0c47aacd905cd12", "7667bf94424f30ad4d2b77d86e9c0cbbc4a9925d", "a29db3ebf677f1fe317ca4ecf0a65a172d4735be", "99746c752c280a60ffb906c7af7d905b86100dab", "fa7c7ef9d487341f806fc5a150486e14565f4706", "8bce07cc60bba1aff1eb250701cd67eb7b349f3f", "7566c26e4b0e51d84900b8f153fc6f069ad09ef7",
				"21c90d16875689f7f8f2ab74158df02dcc340c4d", "ead2a1cbe71522fbc97e7c8ceee4b82ff5ca678e", "0c1f43f2b478261c7bdf92b48e37f318eab78343", "5e59c434b5633b0c3b20062c312e58935724ae2e", "5909ff98ca17534d1d99e7f1cbbdc4f380038016", "6b51154f8be5ac72e5348a17f358e66ce59f8e29", "5f869581d7d6444381310b45c211196b4348c985",
				"a663899ab9369fec8efecfde4d11618a9f1f9467", "5b0b41c364c89c5cb13f011823e0d6ee9b89af26", "a87ee0c50662b13f11b6fdd3eefd4cee17930599", "c13100a1f96e8f512259eda0749ef40ac707e54b", "e2b5ede89ad9d9b811e975f03fb3d7b69b2c90fa", "4841ccaac983b40eff8c7724afd31f4163277cbe", "e74c3159aaace29a29c7ccd722b2d8ecd6815a83",
				"a776378cfe529ea86020de64b67ce6d4aeb9a6a8", "5977336e7fb61c5bee22a2447a36b8455f8ef1a5", "b0f23d36ed6cc1ec327a4200da6f352c972d04ea", "57b83532523bbf553e85fa13e1c90ce389c2dfa1", "17ab18dd54cec0a87bc0d27696306f6695890fdb", "8e316448fc6bfcbef657930a1eda64be39e5eb69", "4b02df3aa8192bf7dab499edfcf8b291ba8dde1f",
				"afe4d367f9f9634b6a993467dd19a7baf7b48f7e", "47447af13e4bc96c143a500034a9182cd569ad07", "f0da9907c8c1466f57c6114997719929b2b96239", "80d83a515872832dff87986f3791a646f52ec1e0", "59b359ed72e979a0560851a6f6acc61ab3f5ec9e", "ba98d13900848f9b3d62a56074129fbdee25300e", "a71bfa04b5a9982cb764b1176d8bd44777247a88",
				"9f3f751faa031e95e7e2b5269a9e3f2189537bde", "8b62ac801b7627562beb06bf6134e1972eafad2a", "d653f67638d9f9e7bd0f65e75c3b349554e016b8", "6a99cef1b67c7f54ec6b3d43d2602972baf7416b", "f1dff225e758dd5bc8ab8b91894096215297b2be", "93d669e435e44d85d5c52b6da70c6ce4ea42eb35", "60a3fae23e73c566fbd6cad622a74aa6977eec5f",
				"54808df261e5b67f6f9b89af910a894dfaf3aa25", "a157cb04b654955eb4b96d91ac392ed71b77ab3d", "250650fa496eff51c2a6bb42e71a63bb285252bf", "d6f5521d37c24f55967e9a15aea8ce8609b749c7", "58ff572db7c8dfbb30d0dbb81717b9641c453a58", "7fa2caefeabb57cf090a2fc3e6eda3a3ce1b8a55", "3ddc1090db671fe0ef11769a2e21e5557c0f865a",
				"9c1f46a30fffcf614d2ec57ce7857a829bbe8878", "5af95fa4195d82b4207df9e5fc2b31f5656be6fc", "b8be0e93b423dad221abe29545fbe8ec36e806bc", "f77a89cd31769d34731ad06cd0c965557826dc25", "f3e79bf8229736a9f3ae208725574436e9d4ac03", "3b1f439733ab2b517bd7065f28b8c888bbbfd170", "0f1e4e91598af7d14ca447b9563e7730a8dd36ac",
				"8b01f8fd69743b0cfa4b611bf8d58a4a60c46183", "b0a25c726258a200f6e8588c402d6bb3d756347e", "13bdbea145163d779aa44cf3cdc8c666348e9b45", "cdd9b50963b6da5cc28a6e71eb529044509f58a8", "2f06e50a17e0425b72826bc4827472b694ed46fe", "d7e3e5d40da7ad5df884115ca8d2876dcd45b956", "6caae01aaa534cc4cb2cb1d8d0a8fd4a9553b155",
				"9d3557b15da0ce44aee7d0dfe1d955139d6da4b9", "0017fb2ac706d21d2299871b988a121b6d6798b7", "01c7bfdd1896629293f02b7ac5b3c4670b78c0ac", "96043a243e14adf367281cc9e8b6a38b554f4725", "3df8bbe8d35bb91cc3e465eebd2f8f5e9975e129", "2399a3276d765d47adbb08417cf84d1d9a6c3e3a", "32b5dc9f00d5ca82080320ae945d4aa2f3a8722a",
				"e3ffc92814fc4b618630521091642afd0ca5a4dc", "df06f065c02315ad7e3849cde3aca5d2bebd2e69", "695bf734b23fccc50c4a21d9f450c62035b93a16", "4312dd61199337c176316f7f1e4a7d1388d0d015", "39065ea36ccf9789327aab73ea88f182c8b77bd3", "226ab3eec3d2075a2d2ba63c64f9bb141ea61be1", "aac4cb4d0d1b1f876eb2dfd2e83ef0f264038168",
				"167db7973bd489807bafec53f5ec283c637a5328", "fd0594a6714a69bdec951ab005c20e0bf46901ff", "52761d617bfc1280cddf02d62aa6ac29fbd98d61", "a8cca58a9a49b08a1325be5fe76646ea85201dbd", "242a2615aab8d88a939ffb5f16e01d2e862f47cf", "0d22988fc5342e7db1ad344c23c89fd239fe46b0", "bef0ccea3da2cbaf92814999e2b933797d6a820c",
				"e044973cad9611dde67e9a21eebf7654613c522a", "369cbc850685e373d936dfdb773e1e71049ca7a2", "95ff1db7f10825325d806a1e336f430a78bde672", "a91ec93a8ffb2e0955dc940297f49f668056f587", "ac9f4b88f3cf4897bb39af2ac4f8f81c09583d4b", "abf8c3b0426cd89fce01770a569d525c648a92b5", "d2af9d4cec08f060dfe47510f6b709ebf01d5686",
				"998ff1746837e487123acd4d17915d1214c9db5b", "8af62859ffc14987a76f128664ccfa0311e0df1f", "f8d42f0b663163bcefe71bd020053e43fab2b3e3", "49fe7ca8bb2f7a794b1cba1d877d98dae520ac73", "5fcb9eca22946e79fe9a2cde6fcf536c4e970a28", "d719ce0fca4016f28cfefbc950bf5d85b11fd430", "3883c0e0590d33b0de17683864a219169ce42e41",
				"a30c5c8cb8faf78b7ad82118e8fd7039901a5a32", "18b0d45be9f72c3c20a47f992325cb0f8af0fe7c", "65200fad3d2b6b48666852efda6f5baa7010fd61", "32d4c07ecf01f2aeee3f07f9b170d9798b5e1d37", "af66d7d133ecbb9f51955fd6591959289a101c01", "8aeaa89d4b512e8be9c377282146833763dcb685", "5a7bc5cc4ff1bf4410336e161fa02379d17f9385",
				"f3db5351ece03ff169e23a785f43b60a3a0315bf", "f6ec70315c2681a3bad9a0eef1017478f6ebed3c", "bc614074d3d37b953898936988dc6744de9a7633", "e19a4aebf5bef7a76ad8560ec51af0388899f6a1", "f72a74d78fca219687c4939d642680d583ce5667", "6adbb83be0f070d0581d585f947a96462214c49a", "9392b51a56555dccff92288b49590e2b690272c9",
				"4b9b7aa6e92a8111b1d6274482a3c774699e47ad", "97d43997cd39161e702811d416da64599a7eea9f", "05ad3ee77cf23214a19e70c391fc8c8fc144a355", "ce022e39c8e4a852fbc77d0db14643b1d0056ad5", "c248ec227ab0aa80760717ed76a0b38128d41e41", "21fc790595a1b79d59087fabaab1a7fe9fde025b", "53d1950fb08d3a1dd5efc790ca869dcd1a32dc5b",
				"efe149c2b6c571c585e487e067c1eddbfd5e8f32", "267953cff8ffaf7f43e5e92b48927c9fc03c7a40", "a19f881b29ba4080038cf8db6d36f6199e89d1ef", "2b6c71a57bad61e17276ba9b9e4c58959cad1d7b", "eb31fab1fa2a32bf3b21068fd75d6b1fdeda9479", "e4264fc7a7911ce8083a61028fe47e49c74100cf", "5a814bd051c02169b2463a437ae35cb90dc9cb83",
				"8b3319dde25b45223a24f909d98303cbf2c878f1", "7913d4b05f7515a5e201d9e29ebb7f2c6a722cce", "f13b26bb384035f77a717f63b0a8afe9f99f30bb", "1fa91fafc5dd2caba9762b13756084626344f220", "9d4c31875f534dc1f989db6c13b92ab1aab85ecf", "719b2540e63955fb6a90bc4f2c4fd9cfd8724e1a", "e49d0548de2488813fa4755adb825cfce4535539",
				"7317cdf6f62be93b22295062e191f6ba59a5db26",
			},
			"not_channel_ids":          []string{},
			"order_by":                 []string{"release_time"},
			"fee_amount":               "<=0",
			"release_time":             ">" + fmt.Sprint(time.Now().Unix()-15778458),
			"include_purchase_receipt": true,
		},
	}
	claimSearchReqData, _ := json.Marshal(claimSearchData)
	frontpageDataRes, err := http.Post(viper.GetString("API_URL")+"/api/v1/proxy?m=claim_search", "application/json", bytes.NewBuffer(claimSearchReqData))
	if err != nil {
		fmt.Println(err)
	}

	frontpageDataBody, err := ioutil.ReadAll(frontpageDataRes.Body)
	if err != nil {
		fmt.Println(err)
	}

	videos := make([]types.Video, 0)
	videosData := gjson.Parse(string(frontpageDataBody))

	waitingVideos.Add(int(videosData.Get("result.items.#").Int()))
	videosData.Get("result.items").ForEach(
		func(key gjson.Result, value gjson.Result) bool {
			go func() {
				claimId := value.Get("claim_id").String()
				lbryUrl := value.Get("canonical_url").String()
				channelLbryUrl := value.Get("signing_channel.canonical_url").String()

				time := time.Unix(value.Get("value.release_time").Int(), 0)
				thumbnail := value.Get("value.thumbnail.url").String()

				videos = append(videos, types.Video{
					Url:       utils.LbryTo(lbryUrl, "http"),
					LbryUrl:   lbryUrl,
					RelUrl:    utils.LbryTo(lbryUrl, "rel"),
					OdyseeUrl: utils.LbryTo(lbryUrl, "odysee"),
					ClaimId:   value.Get("claim_id").String(),
					Channel: types.Channel{
						Name:      value.Get("signing_channel.name").String(),
						Title:     value.Get("signing_channel.value.title").String(),
						Id:        value.Get("signing_channel.claim_id").String(),
						Url:       utils.LbryTo(channelLbryUrl, "http"),
						RelUrl:    utils.LbryTo(channelLbryUrl, "rel"),
						OdyseeUrl: utils.LbryTo(channelLbryUrl, "odysee"),
					},
					Title:        value.Get("value.title").String(),
					ThumbnailUrl: "/image?url=" + thumbnail + "&hash=" + utils.EncodeHMAC(thumbnail),
					Views:        GetVideoViews(claimId),
					Timestamp:    time.Unix(),
					Date:         time.Month().String() + " " + fmt.Sprint(time.Day()) + ", " + fmt.Sprint(time.Year()),
					Duration:     utils.FormatDuration(value.Get("value.video.duration").Int()),
					RelTime:      humanize.Time(time),
				})
				waitingVideos.Done()
			}()

			return true
		},
	)
	waitingVideos.Wait()

	commentCache.Set("fp", videos, cache.DefaultExpiration)
	return videos
}
