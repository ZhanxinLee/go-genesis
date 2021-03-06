// MIT License
//
// Copyright (c) 2016-2018 GenesisKernel
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package api

import (
	"fmt"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/GenesisKernel/go-genesis/packages/conf"
	"github.com/GenesisKernel/go-genesis/packages/consts"
	"github.com/GenesisKernel/go-genesis/packages/converter"
	"github.com/GenesisKernel/go-genesis/packages/crypto"
	taskContract "github.com/GenesisKernel/go-genesis/packages/scheduler/contract"
)

func TestVDECreate(t *testing.T) {
	var (
		err   error
		retid int64
		ret   vdeCreateResult
	)
	if err = keyLogin(1); err != nil {
		t.Error(err)
		return
	}

	if err = sendPost(`vde/create`, nil, &ret); err != nil &&
		err.Error() != `400 {"error": "E_VDECREATED", "msg": "Virtual Dedicated Ecosystem is already created" }` {
		t.Error(err)
		return
	}

	rnd := `rnd` + crypto.RandSeq(6)
	form := url.Values{`Value`: {`contract ` + rnd + ` {
		    data {
				Par string
			}
			action { Test("active",  $Par)}}`}, `Conditions`: {`ContractConditions("MainCondition")`}, `vde`: {`true`}}

	if retid, _, err = postTxResult(`NewContract`, &form); err != nil {
		t.Error(err)
		return
	}
	form = url.Values{`Id`: {converter.Int64ToStr(retid)}, `Value`: {`contract ` + rnd + ` {
		data {
			Par string
		}
		action { Test("active 5",  $Par)}}`}, `Conditions`: {`ContractConditions("MainCondition")`}, `vde`: {`true`}}

	if err := postTx(`EditContract`, &form); err != nil {
		t.Error(err)
		return
	}
	form = url.Values{`Name`: {rnd}, `Value`: {`Test value`}, `Conditions`: {`ContractConditions("MainCondition")`},
		`vde`: {`1`}}

	if retid, _, err = postTxResult(`NewParameter`, &form); err != nil {
		t.Error(err)
		return
	}
	form = url.Values{`Name`: {`new_table`}, `Value`: {`Test value`}, `Conditions`: {`ContractConditions("MainCondition")`},
		`vde`: {`1`}}
	if err = postTx(`NewParameter`, &form); err != nil && err.Error() !=
		`500 {"error": "E_SERVER", "msg": "{\"type\":\"warning\",\"error\":\"Parameter new_table already exists\"}" }` {
		t.Error(err)
		return
	}
	form = url.Values{`Id`: {converter.Int64ToStr(retid)}, `Value`: {`Test edit value`}, `Conditions`: {`true`},
		`vde`: {`1`}}
	if _, _, err = postTxResult(`EditParameter`, &form); err != nil {
		t.Error(err)
		return
	}

	form = url.Values{"Name": {`menu` + rnd}, "Value": {`first
		second
		third`}, "Title": {`My Menu`},
		"Conditions": {`true`}, `vde`: {`1`}}
	retid, _, err = postTxResult(`NewMenu`, &form)
	if err != nil {
		t.Error(err)
		return
	}
	form = url.Values{`Id`: {converter.Int64ToStr(retid)}, `Value`: {`Test edit value`},
		`Conditions`: {`true`},
		`vde`:        {`1`}}
	if err = postTx(`EditMenu`, &form); err != nil {
		t.Error(err)
		return
	}
	form = url.Values{"Id": {converter.Int64ToStr(retid)}, "Value": {`Span(Append)`},
		`vde`: {`1`}}
	err = postTx(`AppendMenu`, &form)
	if err != nil {
		t.Error(err)
		return
	}

	form = url.Values{"Name": {`page` + rnd}, "Value": {`Page`}, "Menu": {`government`},
		"Conditions": {`true`}, `vde`: {`1`}}
	retid, _, err = postTxResult(`NewPage`, &form)
	if err != nil {
		t.Error(err)
		return
	}
	form = url.Values{`Id`: {converter.Int64ToStr(retid)}, `Value`: {`Test edit page value`},
		`Conditions`: {`true`}, "Menu": {`government`},
		`vde`: {`1`}}
	if err = postTx(`EditPage`, &form); err != nil {
		t.Error(err)
		return
	}
	form = url.Values{"Id": {converter.Int64ToStr(retid)}, "Value": {`Span(Test Page)`},
		`vde`: {`1`}}
	err = postTx(`AppendPage`, &form)
	if err != nil {
		t.Error(err)
		return
	}
	form = url.Values{"Name": {`block` + rnd}, "Value": {`Page block`}, "Conditions": {`true`}, `vde`: {`1`}}
	retid, _, err = postTxResult(`NewBlock`, &form)
	if err != nil {
		t.Error(err)
		return
	}
	form = url.Values{`Id`: {converter.Int64ToStr(retid)}, `Value`: {`Test edit block value`},
		`Conditions`: {`true`}, `vde`: {`1`}}
	if err = postTx(`EditBlock`, &form); err != nil {
		t.Error(err)
		return
	}

	name := randName(`tbl`)
	form = url.Values{"Name": {name}, `vde`: {`true`}, "Columns": {`[{"name":"MyName","type":"varchar", "index": "1",
			  "conditions":"true"},
			{"name":"Amount", "type":"number","index": "0", "conditions":"true"},
			{"name":"Active", "type":"character","index": "0", "conditions":"true"}]`},
		"Permissions": {`{"insert": "true", "update" : "true", "new_column": "true"}`}}
	err = postTx(`NewTable`, &form)
	if err != nil {
		t.Error(err)
		return
	}
	form = url.Values{"Name": {name}, `vde`: {`true`},
		"Permissions": {`{"insert": "ContractConditions(\"MainCondition\")",
						"update" : "true", "new_column": "ContractConditions(\"MainCondition\")"}`}}
	err = postTx(`EditTable`, &form)
	if err != nil {
		t.Error(err)
		return
	}
	form = url.Values{"TableName": {name}, "Name": {`newCol`}, `vde`: {`1`},
		"Type": {"varchar"}, "Index": {"0"}, "Permissions": {"true"}}
	err = postTx(`NewColumn`, &form)
	if err != nil {
		t.Error(err)
		return
	}
	form = url.Values{"TableName": {name}, "Name": {`newColRead`}, `vde`: {`1`},
		"Type": {"varchar"}, "Index": {"0"}, "Permissions": {`{"update":"true", "read":"false"}`}}
	err = postTx(`NewColumn`, &form)
	if err != nil {
		t.Error(err)
		return
	}

	form = url.Values{"TableName": {name}, "Name": {`newCol`}, `vde`: {`1`},
		"Permissions": {"ContractConditions(\"MainCondition\")"}}
	err = postTx(`EditColumn`, &form)
	if err != nil {
		t.Error(err)
		return
	}
	form = url.Values{"TableName": {name}, "Name": {`newCol`}, `vde`: {`1`},
		"Permissions": {`{"update":"true", "read":"false"}`}}
	err = postTx(`EditColumn`, &form)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestVDEParams(t *testing.T) {
	if err := keyLogin(1); err != nil {
		t.Error(err)
		return
	}
	rnd := `rnd` + crypto.RandSeq(6)
	form := url.Values{`Name`: {rnd}, `Value`: {`Test value`}, `Conditions`: {`ContractConditions("MainCondition")`},
		`vde`: {`true`}}
	if _, _, err := postTxResult(`NewParameter`, &form); err != nil {
		t.Error(err)
		return
	}

	var ret ecosystemParamsResult
	err := sendGet(`ecosystemparams?vde=true`, nil, &ret)
	if err != nil {
		t.Error(err)
		return
	}
	if len(ret.List) < 5 {
		t.Error(fmt.Errorf(`wrong count of parameters %d`, len(ret.List)))
	}
	err = sendGet(`ecosystemparams?vde=true&names=stylesheet,`+rnd, nil, &ret)
	if err != nil {
		t.Error(err)
		return
	}
	if len(ret.List) != 2 {
		t.Error(fmt.Errorf(`wrong count of parameters %d`, len(ret.List)))
	}
	var parValue paramValue
	err = sendGet(`ecosystemparam/`+rnd+`?vde=true`, nil, &parValue)
	if err != nil {
		t.Error(err)
		return
	}
	if parValue.Name != rnd {
		t.Error(fmt.Errorf(`wrong value of parameter`))
	}
	var tblResult tablesResult
	err = sendGet(`tables?vde=true`, nil, &tblResult)
	if err != nil {
		t.Error(err)
		return
	}
	if tblResult.Count < 5 {
		t.Error(fmt.Errorf(`wrong tables result`))
	}
	form = url.Values{"Name": {rnd}, `vde`: {`1`}, "Columns": {`[{"name":"MyName","type":"varchar", "index": "1",
		"conditions":"true"},
	  {"name":"Amount", "type":"number","index": "0", "conditions":"true"},
	  {"name":"Active", "type":"character","index": "0", "conditions":"true"}]`},
		"Permissions": {`{"insert": "true", "update" : "true", "new_column": "true"}`}}
	err = postTx(`NewTable`, &form)
	if err != nil {
		t.Error(err)
		return
	}
	var tResult tableResult
	err = sendGet(`table/`+rnd+`?vde=true`, nil, &tResult)
	if err != nil {
		t.Error(err)
		return
	}
	if tResult.Name != rnd {
		t.Error(fmt.Errorf(`wrong table result`))
		return
	}
	var retList listResult
	err = sendGet(`list/contracts?vde=true`, nil, &retList)
	if err != nil {
		t.Error(err)
		return
	}
	if converter.StrToInt64(retList.Count) < 7 {
		t.Error(fmt.Errorf(`The number of records %s < 7`, retList.Count))
		return
	}
	var retRow rowResult
	err = sendGet(`row/contracts/2?vde=true`, nil, &retRow)
	if err != nil {
		t.Error(err)
		return
	}
	if !strings.Contains(retRow.Value[`value`], `VDEFunctions`) {
		t.Error(`wrong row result`)
		return
	}
	var retCont contractsResult
	err = sendGet(`contracts?vde=true`, nil, &retCont)
	if err != nil {
		t.Error(err)
		return
	}
	form = url.Values{`Value`: {`contract ` + rnd + ` {
		data {
			Par string
		}
		action { Test("active",  $Par)}}`}, `Conditions`: {`ContractConditions("MainCondition")`}, `vde`: {`true`}}

	if _, _, err = postTxResult(`NewContract`, &form); err != nil {
		t.Error(err)
		return
	}
	var cont getContractResult
	err = sendGet(`contract/`+rnd+`?vde=true`, nil, &cont)
	if err != nil {
		t.Error(err)
		return
	}
	if !strings.HasSuffix(cont.Name, rnd) {
		t.Error(`wrong contract result`)
		return
	}

	form = url.Values{"Name": {rnd}, "Value": {`Page`}, "Menu": {`government`},
		"Conditions": {`true`}, `vde`: {`1`}}
	err = postTx(`NewPage`, &form)
	if err != nil {
		t.Error(err)
		return
	}
	err = sendPost(`content/page/`+rnd, &url.Values{`vde`: {`true`}}, &ret)
	if err != nil {
		t.Error(err)
		return
	}
	form = url.Values{"Name": {rnd}, "Value": {`Menu`}, "Conditions": {`true`}, `vde`: {`1`}}
	err = postTx(`NewMenu`, &form)
	if err != nil {
		t.Error(err)
		return
	}
	err = sendPost(`content/menu/`+rnd, &url.Values{`vde`: {`true`}}, &ret)
	if err != nil {
		t.Error(err)
		return
	}

	name := randName(`lng`)
	value := `{"en": "My VDE test", "fr": "French VDE test"}`

	form = url.Values{"Name": {name}, "Trans": {value}, "vde": {`true`}}
	err = postTx(`NewLang`, &form)
	if err != nil {
		t.Error(err)
		return
	}
	input := fmt.Sprintf(`Span($%s$)+LangRes(%[1]s,fr)`, name)
	var retContent contentResult
	err = sendPost(`content`, &url.Values{`template`: {input}, `vde`: {`true`}}, &retContent)
	if err != nil {
		t.Error(err)
		return
	}
	if RawToString(retContent.Tree) != `[{"tag":"span","children":[{"tag":"text","text":"My VDE test"}]},{"tag":"text","text":"+French VDE test"}]` {
		t.Error(fmt.Errorf(`wrong tree %s`, RawToString(retContent.Tree)))
		return
	}

	name = crypto.RandSeq(4)
	err = postTx(`Import`, &url.Values{"vde": {`true`}, "Data": {fmt.Sprintf(imp, name)}})
	if err != nil {
		t.Error(err)
		return
	}
}

var vdeimp = `{
    "pages": [
        {
            "Name": "imp_page2",
            "Conditions": "true",
            "Menu": "imp",
            "Value": "imp"
        }
    ],
    "blocks": [
        {
            "Name": "imp2",
            "Conditions": "true",
            "Value": "imp"
        }
    ],
    "menus": [
        {
            "Name": "imp2",
            "Conditions": "true",
            "Value": "imp"
        }
    ],
    "parameters": [
        {
            "Name": "founder_account2",
            "Value": "-6457397116804798941",
            "Conditions": "ContractConditions(\"MainCondition\")"
        },
        {
            "Name": "test_pa2",
            "Value": "1",
            "Conditions": "true"
        }
    ],
    "languages": [
        {
            "Name": "est2",
            "Trans": "{\"en\":\"yeye\",\"te\":\"knfek\"}"
        }
    ],
    "contracts": [
        {
            "Name": "testCont2",
            "Value": "contract testCont2 {\n    data {\n\n    }\n\n    conditions {\n\n    }\n\n    action {\n        $result=\"privet\"\n    }\n}",
            "Conditions": "true"
        }
    ],
    "tables": [
        {
            "Name": "tests2",
            "Columns": "[{\"name\":\"name\",\"type\":\"text\",\"conditions\":\"true\"}]",
            "Permissions": "{\"insert\":\"true\",\"update\":\"true\",\"new_column\":\"true\"}"
        }
    ],
    "data": [
        {
            "Table": "tests2",
            "Columns": [
                "name"
            ],
            "Data": []
        }
    ]
}`

func TestVDEImport(t *testing.T) {
	if err := keyLogin(1); err != nil {
		t.Error(err)
		return
	}
	err := postTx(`Import`, &url.Values{"vde": {`true`}, "Data": {vdeimp}})
	if err != nil {
		t.Error(err)
		return
	}
}

func TestHTTPRequest(t *testing.T) {
	if err := keyLogin(1); err != nil {
		t.Error(err)
		return
	}
	rnd := `rnd` + crypto.RandSeq(6)
	form := url.Values{`Value`: {`contract ` + rnd + ` {
		    data {
				Auth string
			}
			action {
				var ret string 
				var pars, heads, json map
				ret = HTTPRequest("http://www.instagram.com/", "GET", heads, pars)
				if !Contains(ret, "react-root") {
					error "instagram error"
				}
				ret = HTTPRequest("http://www.google.com/search?q=exotic", "GET", heads, pars)
				if !Contains(ret, "exotic") {
					error "google error"
				}
				heads["Authorization"] = "Bearer " + $Auth
				pars["vde"] = "true"
				ret = HTTPRequest("http://localhost:7079` + consts.ApiPath + `content/page/` + rnd + `", "POST", heads, pars)
				json = JSONToMap(ret)
				if json["menu"] != "myvdemenu" {
					error "Wrong vde menu"
				}
				ret = HTTPRequest("http://localhost:7079` + consts.ApiPath + `contract/VDEFunctions?vde=true", "GET", heads, pars)
				json = JSONToMap(ret)
				if json["name"] != "@1VDEFunctions" {
					error "Wrong vde contract"
				}
			}}`}, `Conditions`: {`true`}, `vde`: {`true`}}

	if err := postTx(`NewContract`, &form); err != nil {
		t.Error(err)
		return
	}
	form = url.Values{"Name": {rnd}, "Value": {`Page`}, "Menu": {`myvdemenu`},
		"Conditions": {`true`}, `vde`: {`true`}}
	if err := postTx(`NewPage`, &form); err != nil {
		t.Error(err)
		return
	}
	if err := postTx(rnd, &url.Values{`vde`: {`true`}, `Auth`: {gAuth}}); err != nil {
		t.Error(err)
		return
	}
}

func TestNodeHTTPRequest(t *testing.T) {
	var err error
	if err = keyLogin(1); err != nil {
		t.Error(err)
		return
	}

	rnd := `rnd` + crypto.RandSeq(4)

	form := url.Values{`Value`: {`contract for` + rnd + ` {
		data {
			Par string
		}
		action { $result = "Test NodeContract " + $Par + " ` + rnd + `"}
    }`}, `Conditions`: {`ContractConditions("MainCondition")`}}

	if err = postTx(`NewContract`, &form); err != nil {
		t.Error(err)
		return
	}
	var ret getContractResult
	err = sendGet(`contract/for`+rnd, nil, &ret)
	if err != nil {
		t.Error(err)
		return
	}
	if err := postTx(`ActivateContract`, &url.Values{`Id`: {ret.TableID}}); err != nil {
		t.Error(err)
		return
	}

	form = url.Values{`Value`: {`contract ` + rnd + ` {
		    data {
				Par string
			}
			action {
				var ret string 
				var pars, heads, json map
				heads["Authorization"] = "Bearer " + $auth_token
				pars["vde"] = "false"
				pars["Par"] = $Par
				ret = HTTPRequest("http://localhost:7079` + consts.ApiPath + `node/for` + rnd + `", "POST", heads, pars)
				json = JSONToMap(ret)
				$result = json["hash"]
			}}`}, `Conditions`: {`true`}, `vde`: {`true`}}

	if err = postTx(`NewContract`, &form); err != nil {
		t.Error(err)
		return
	}
	var (
		msg string
		id  int64
	)
	if _, msg, err = postTxResult(rnd, &url.Values{`vde`: {`true`}, `Par`: {`node`}}); err != nil {
		t.Error(err)
		return
	}
	id, err = waitTx(msg)
	if id != 0 && err != nil {
		msg = err.Error()
		err = nil
	}
	if msg != `Test NodeContract node `+rnd {
		t.Error(`wrong result: ` + msg)
	}
	form = url.Values{`Value`: {`contract node` + rnd + ` {
		data {
		}
		action { 
			var ret string 
			var pars, heads, json map
			heads["Authorization"] = "Bearer " + $auth_token
			pars["vde"] = "false"
			pars["Par"] = "NodeContract testing"
			ret = HTTPRequest("http://localhost:7079` + consts.ApiPath + `node/for` + rnd + `", "POST", heads, pars)
			json = JSONToMap(ret)
			$result = json["hash"]
		}
    }`}, `Conditions`: {`ContractConditions("MainCondition")`}, `vde`: {`true`}}

	if err = postTx(`NewContract`, &form); err != nil {
		t.Error(err)
		return
	}
	// You can specify the directory with NodePrivateKey & NodePublicKey files
	conf.Config.PrivateDir = ``
	if len(conf.Config.PrivateDir) > 0 {
		conf.Config.HTTP.Host = `localhost`
		conf.Config.HTTP.Port = 7079

		nodeResult, err := taskContract.NodeContract(`@1node` + rnd)
		if err != nil {
			t.Error(err)
			return
		}
		id, err = waitTx(nodeResult.Result)
		if id != 0 && err != nil {
			msg = err.Error()
			err = nil
		}
		if msg != `Test NodeContract NodeContract testing `+rnd {
			t.Error(`wrong result: ` + msg)
		}
	}
}

func TestCron(t *testing.T) {
	if err := keyLogin(1); err != nil {
		t.Error(err)
		return
	}

	err := postTx("NewCron", &url.Values{
		"Cron":       {"60 * * * * *"},
		"Contract":   {"TestCron"},
		"Conditions": {`ContractConditions("MainCondition")`},
		"vde":        {"true"},
	})
	if err.Error() != `500 {"error": "E_SERVER", "msg": "{\"type\":\"panic\",\"error\":\"End of range (60) above maximum (59): 60\"}" }` {
		t.Error(err)
	}

	postTx("NewContract", &url.Values{
		"Value": {`
			contract TestCron {
				data {}
				action {
					return "Success"
				}
			}
		`},
		"Conditions": {`ContractConditions("MainCondition")`},
		"vde":        {"true"},
	})

	till := time.Now().Format(time.RFC3339)
	err = postTx("NewCron", &url.Values{
		"Cron":       {"* * * * * *"},
		"Contract":   {"TestCron"},
		"Conditions": {`ContractConditions("MainCondition")`},
		"Till":       {till},
		"vde":        {"true"},
	})
	if err != nil {
		t.Error(err)
	}

	err = postTx("EditCron", &url.Values{
		"Id":         {"1"},
		"Cron":       {"*/3 * * * * *"},
		"Contract":   {"TestCron"},
		"Conditions": {`ContractConditions("MainCondition")`},
		"Till":       {till},
		"vde":        {"true"},
	})
	if err != nil {
		t.Error(err)
	}
}
