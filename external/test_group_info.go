package external

import (
	"coresamples/common"
	_ "embed"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt"
	capi "github.com/hashicorp/consul/api"
	"io"
	"net/http"
	"time"
)

////go:embed test_group.json
//var testGroup []byte
//
////go:embed test_group_mapping.json
//var testGroupMapping []byte

type SpecialProduct int

/**
/**
  LIST_OF_ALL_IBSSURE_REPORT_GROUPS = ReportVersion19.LIST_OF_ALL_IBSSURE_REPORT_GROUPS;
  LIST_OF_ALL_MICRNONUTRIENTS_GROUPS = ReportVersion20.LIST_OF_ALL_MICRNONUTRIENTS_GROUPS; //
    WHEAT_ZOOMER_GROUPS = ReportVersion20.WHEAT_ZOOMER_GROUPS; //
    LIST_OF_ALL_FOOD_SENSITIVITY_REPORT_GROUPS = ReportVersion20.LIST_OF_ALL_FOOD_SENSITIVITY_REPORT_GROUPS; //
  GUT_BACTERIA_GROUPS_V2 = reportVersion == ReportVersion18.GUT_BACTERIA_GROUPS_V2;
  WELLNESS_CARDIAX_GROUPS_ALL = ReportVersion15.WELLNESS_CARDIAX_V1_GROUPS; // old version
  LIST_OF_ALL_RESPIRATORY_VIRUS_REPORT_GROUPS = ReportVersion17.LIST_OF_ALL_RESPIRATORY_VIRUS_REPORT_GROUPS;
  LIST_OF_ALL_VG_COMMENSAL_BACTERIA_GROUPS = ReportVersion20.LIST_OF_ALL_VG_COMMENSAL_BACTERIA_GROUPS;
*/

const (
	TypeGroup                       = "GROUP"
	TypeTest                        = "TEST"
	GutBacteriaGroup SpecialProduct = iota
	FoodSensitivityReportGroup
	WheatZoomerGroup
	MicroNeutrientsGroup
	IBSSureReportGroup
	RespiratoryVirusReportGroup
	VGCommensalBacteriaGroup
	CardiaxGroup
)

type TestGroup struct {
	Id          int    `json:"id,omitempty"`
	OrderTypeId int    `json:"orderTypeId,omitempty"` // order type ID, not unique, can be shared by group and test
	GroupType   string `json:"type"`                  // whether a test group is a group or a single test
	Name        string `json:"name"`
}

type TestGroupInfo struct {
	Groups            map[int]*TestGroup
	GroupTestMappings map[int][]int // group order type ID -> test ID

	// tests in special product groups, stored in set
	SpecialGroupTests map[SpecialProduct]map[int]bool
}

func (i *TestGroupInfo) ReSync(client *capi.Client, prefix string) {
	wait := time.Second

	if common.Env.RunEnv == common.AksProductionEnv {
		common.InitEndpointsFromConsul(client, common.Env.ConsulPrefix, "endpoints")
	} else {
		common.InitEndpointsFromConsul(client, common.Env.ConsulPrefix, "endpointsStaging")
	}
	if common.EndpointsInfo != nil && common.EndpointsInfo.GetTestGroupMapping != "" && common.EndpointsInfo.GetTestGroup != "" {
		i.Init()
	}

	for len(i.GroupTestMappings) == 0 || len(i.SpecialGroupTests) == 0 || len(i.Groups) == 0 {
		if wait > time.Second*10 {
			break
		}

		time.Sleep(wait)
		wait = wait * 2

		if common.Env.RunEnv == common.AksProductionEnv {
			common.InitEndpointsFromConsul(client, common.Env.ConsulPrefix, "endpoints")
		} else {
			common.InitEndpointsFromConsul(client, common.Env.ConsulPrefix, "endpointsStaging")
		}
		if common.EndpointsInfo != nil && common.EndpointsInfo.GetTestGroupMapping != "" && common.EndpointsInfo.GetTestGroup != "" {
			i.Init()
		}
	}
}

func (i *TestGroupInfo) Init() {
	i.getTestGroups()
	i.getGroupTestMapping()
	i.initSpecialGroupTests()
}

func (i *TestGroupInfo) GetTestsUnderGroup(orderTypeId int) []int {
	if ret, ok := i.GroupTestMappings[orderTypeId]; ok {
		return ret
	}
	return nil
}

func (i *TestGroupInfo) GetTestsUnderGroups(orderTypeIds ...int) map[int]bool {
	tests := map[int]bool{}
	for _, id := range orderTypeIds {
		testsUnderGroup := i.GetTestsUnderGroup(id)
		if testsUnderGroup == nil {
			continue
		}
		for _, testId := range testsUnderGroup {
			tests[testId] = true
			break
		}
	}
	return tests
}

func (i *TestGroupInfo) getTestGroups() {
	token, err := i.GenerateToken()
	if err != nil {
		common.Error(err)
	}
	req, _ := http.NewRequest("GET", common.EndpointsInfo.GetTestGroup, nil)
	req.Header.Add("Authorization", token)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		common.Errorf("unable to fetch test group info", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		common.Error(fmt.Errorf("response status code %d: %s", resp.StatusCode, common.EndpointsInfo.GetTestGroup))
	}
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		common.Errorf("unable to read body from response", err)
	}

	if err := json.Unmarshal(body, &i.Groups); err != nil {
		common.Errorf("unable to parse body from response", err)
	}

}

func (i *TestGroupInfo) getGroupTestMapping() {
	token, err := i.GenerateToken()
	if err != nil {
		common.Error(err)
	}

	req, _ := http.NewRequest("GET", common.EndpointsInfo.GetTestGroupMapping, nil)
	req.Header.Add("Authorization", token)

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		common.Errorf("unable to fetch test group info", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		common.Error(fmt.Errorf("response status code %d: %s", resp.StatusCode, common.EndpointsInfo.GetTestGroupMapping))
	}
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		common.Errorf("unable to read body from response", err)
	}

	if err := json.Unmarshal(body, &i.GroupTestMappings); err != nil {
		common.Errorf("unable to parse body from response", err)
	}
}

func (i *TestGroupInfo) initSpecialGroupTests() {
	i.SpecialGroupTests = make(map[SpecialProduct]map[int]bool)
	i.SpecialGroupTests[IBSSureReportGroup] = i.GetTestsUnderGroups(155)
	i.SpecialGroupTests[MicroNeutrientsGroup] = i.GetTestsUnderGroups(164, 165, 167, 166, 168, 169)
	i.SpecialGroupTests[WheatZoomerGroup] = i.GetTestsUnderGroups(141, 33, 35, 36, 71, 37, 38, 34, 39, 89, 351, 352, 357)
	i.SpecialGroupTests[FoodSensitivityReportGroup] = i.GetTestsUnderGroups(118, 119, 120, 121, 170, 171, 509, 510, 125, 126, 127, 128, 129, 130, 131, 132, 133)
	i.SpecialGroupTests[GutBacteriaGroup] = i.GetTestsUnderGroups(88, 77, 78, 79, 80, 81, 82, 83, 86, 90, 92, 84, 85)
	i.SpecialGroupTests[CardiaxGroup] = i.GetTestsUnderGroups(50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 65, 66, 67, 68, 69, 70)
	i.SpecialGroupTests[RespiratoryVirusReportGroup] = i.GetTestsUnderGroups(98)
	i.SpecialGroupTests[VGCommensalBacteriaGroup] = i.GetTestsUnderGroups(174, 175, 176, 177, 178, 179, 180, 181, 182)
}

func (s *TestGroupInfo) GenerateToken() (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"iat": now.Unix(),
		"exp": now.AddDate(1, 0, 0).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(common.Secrets.Secret))
	if err != nil {
		return "", err
	}
	return "Bearer " + tokenString, nil
}
