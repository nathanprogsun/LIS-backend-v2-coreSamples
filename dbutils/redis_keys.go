package dbutils

import "fmt"

const (
	ServicePrefix = "lis::core_service_v2"
	TopicLabOrder = "lab_order"
	TopicTest     = "test"
	TopicOrder    = "order_info"
	TopicCustomer = "customer"
)

func KeyLabOrderSendRecord(sampleId int, tubeType string) string {
	return fmt.Sprintf("%s::%s::lab_order_send_record_sample_id%d_tube_type%s",
		ServicePrefix, TopicLabOrder, sampleId, tubeType)
}

func KeyGetSampleTests(sampleID int) string {
	return fmt.Sprintf("%s::%s::get_sample_tests::%d",
		ServicePrefix, TopicTest, sampleID)
}

// Key for retrieving sample tube counts
func KeySampleTubeCounts(sampleID int) string {
	return fmt.Sprintf("%s::get_sample_tube_counts_%d", ServicePrefix, sampleID)
}

// Key for retrieving tube type via sample type code
func KeyTubeTypeViaSampleTypeCode(sampleTypeCode string) string {
	return fmt.Sprintf("%s::getTubeTypeViaSampleTypeCode::sample_type_code_%s",
		ServicePrefix, sampleTypeCode)
}

func KeyLabTestsBySampleID(sampleID int) string {
	return fmt.Sprintf("%s::%s::getLabTestsBySampleID_::%d",
		ServicePrefix, TopicOrder, sampleID)
}

func KeyGetCustomerSalesByCustomerName(customerName string) string {
	return fmt.Sprintf("%s::%s::get_customer_sales_name_::%s",
		ServicePrefix, TopicCustomer, customerName)
}

func KeyGetCustomerSalesByCustomerId(customerId int) string {
	return fmt.Sprintf("%s::%s::get_customer_sales_name_::%d",
		ServicePrefix, TopicCustomer, customerId)
}

func KeyGetSalesCustomer(salesName string, page string, perPage string) string {
	return fmt.Sprintf("%s::%s::getSalesCustomer_%s_%s_%s",
		ServicePrefix, TopicCustomer, salesName, page, perPage)
}

func KeyGetCustomerAllClinics(customerId string) string {
	return fmt.Sprintf("%s::%s::get_customer_all_clinics_%s",
		ServicePrefix, TopicCustomer, customerId)
}

func KeyGetTestAll() string {
	return fmt.Sprintf("%s::%s::get_test_all",
		ServicePrefix, TopicTest)
}

func KeyGetTestByTestId(testId int) string {
	return fmt.Sprintf("%s::%s::get_test_test_id_%d",
		ServicePrefix, TopicTest, testId)
}
