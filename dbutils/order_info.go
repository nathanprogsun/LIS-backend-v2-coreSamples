package dbutils

import (
	"context"
	"coresamples/ent"
	"coresamples/ent/labordersendhistory"
	"coresamples/ent/orderinfo"
	"coresamples/ent/sample"
	"fmt"
)

// GetOrderById return the order and corresponding sample and order flags
func GetOrderById(orderId int, client *ent.Client, ctx context.Context) (*ent.OrderInfo, error) {
	return client.OrderInfo.Query().Where(orderinfo.ID(orderId)).WithOrderFlags().WithSample().Only(ctx)
}

func GetOrdersWithSampleId(sampleId int, client *ent.Client, ctx context.Context) ([]*ent.OrderInfo, error) {
	return client.OrderInfo.Query().Where(orderinfo.HasSampleWith(sample.ID(sampleId))).All(ctx)
}

func CancelOrderById(orderId int, client *ent.Client, ctx context.Context) (*ent.OrderInfo, error) {
	err := client.OrderInfo.Update().Where(orderinfo.ID(orderId)).
		SetIsActive(false).
		SetOrderCanceled(true).
		Exec(ctx)
	if err != nil {
		return nil, err
	}

	return client.OrderInfo.Get(ctx, orderId)
}

func RestoreOrderById(orderId int, client *ent.Client, ctx context.Context) (*ent.OrderInfo, error) {
	err := client.OrderInfo.Update().Where(orderinfo.ID(orderId)).
		SetIsActive(true).
		SetOrderCanceled(false).
		Exec(ctx)
	if err != nil {
		return nil, err
	}

	return client.OrderInfo.Get(ctx, orderId)
}

func UpdateOrderStatus(orderId int, orderStatus string, client *ent.Client, ctx context.Context) error {
	_, err := client.OrderInfo.UpdateOneID(orderId).SetOrderStatus(orderStatus).Save(ctx)
	return err
}

func FindOrderWithSampleId(sampleId int32, client *ent.Client, ctx context.Context) (*ent.OrderInfo, error) {
	return client.OrderInfo.Query().Where(
		orderinfo.HasSampleWith(sample.IDEQ(int(sampleId))),
	).First(ctx)
}

func GetLabOrderSendRecord(sampleId int, tubeType string,
	isLabSpecialOrder bool, isResendBlocked bool,
	client *ent.Client, ctx context.Context) (*ent.LabOrderSendHistory, error) {
	return client.LabOrderSendHistory.Query().Where(
		labordersendhistory.And(
			labordersendhistory.SampleIDEQ(sampleId),
			labordersendhistory.TubeTypeEQ(tubeType),
			labordersendhistory.IsLabSpecialOrder(isLabSpecialOrder),
			labordersendhistory.IsResendBlocked(isResendBlocked)),
	).First(ctx)
}

func CreateLabOrderSendRecord(sampleId int, tubeType string, isRedrawOrder bool, isLabSpecialOrder bool, action string,
	client *ent.Client, ctx context.Context) (*ent.LabOrderSendHistory, error) {
	return client.LabOrderSendHistory.Create().
		SetSampleID(sampleId).
		SetTubeType(tubeType).
		SetIsRedrawOrder(isRedrawOrder).
		SetIsLabSpecialOrder(isLabSpecialOrder).
		SetAction(action).
		Save(ctx)
}

// UpdateOrderStatusBySampleID updates a specific field in the order_info table based on the SampleID.
func UpdateOrderStatusBySampleID(sampleID int32, field string, value interface{}, ctx context.Context, client *ent.Client) error {
	// Find the order associated with the sample ID
	order, err := FindOrderWithSampleId(sampleID, client, ctx)

	if err != nil {
		return err
	}

	// Prepare the update based on the field
	update := client.OrderInfo.UpdateOneID(order.ID)

	switch field {
	case orderinfo.FieldOrderStatus:
		update.SetOrderStatus(value.(string))
	case orderinfo.FieldOrderMajorStatus:
		update.SetOrderMajorStatus(value.(string))
	case orderinfo.FieldOrderKitStatus:
		update.SetOrderKitStatus(value.(string))
	case orderinfo.FieldOrderReportStatus:
		update.SetOrderReportStatus(value.(string))
	case orderinfo.FieldOrderTnpIssueStatus:
		update.SetOrderTnpIssueStatus(value.(string))
	case orderinfo.FieldOrderBillingIssueStatus:
		update.SetOrderBillingIssueStatus(value.(string))
	case orderinfo.FieldOrderMissingInfoIssueStatus:
		update.SetOrderMissingInfoIssueStatus(value.(string))
	case orderinfo.FieldOrderIncompleteQuestionnaireIssueStatus:
		update.SetOrderIncompleteQuestionnaireIssueStatus(value.(string))
	case orderinfo.FieldOrderNyWaiveFormIssueStatus:
		update.SetOrderNyWaiveFormIssueStatus(value.(string))
	case orderinfo.FieldOrderLabIssueStatus:
		update.SetOrderLabIssueStatus(value.(string))
	default:
		return err
	}

	// Save the updated order
	if _, err := update.Save(ctx); err != nil {
		return err
	}

	return nil
}

// UpdateMajorOrderStatusBySampleID updates the order_major_status of the order associated with the given sample ID.
func UpdateMajorOrderStatusBySampleID(sampleID int32, newStatus string, ctx context.Context, client *ent.Client) error {
	// Find the sample record by sample ID
	sampleInfo, err := client.Sample.
		Query().
		Where(sample.IDEQ(int(sampleID))).
		Only(ctx)
	if err != nil {
		return fmt.Errorf("failed to find sample with ID %d: %w", sampleID, err)
	}

	// Check if the sample has an associated order ID
	if sampleInfo.OrderID == 0 {
		return fmt.Errorf("no order associated with sample ID %d", sampleID)
	}

	// Update the order_major_status in the order_info table
	_, err = client.OrderInfo.
		UpdateOneID(sampleInfo.OrderID).
		SetOrderMajorStatus(newStatus).
		Save(ctx)
	if err != nil {
		return fmt.Errorf("failed to update order major status for order ID %d: %w", sampleInfo.OrderID, err)
	}

	return nil
}

func GetOrderByAccessionId(accessionId string, client *ent.Client, ctx context.Context) (*ent.OrderInfo, error) {
	return client.OrderInfo.Query().Where(orderinfo.HasSampleWith(sample.AccessionID(accessionId))).First(ctx)
}

func UpdateOrderKitStatus(orderId int, orderKitStatus string, client *ent.Client, ctx context.Context) error {
	_, err := client.OrderInfo.UpdateOneID(orderId).SetOrderKitStatus(orderKitStatus).Save(ctx)
	return err
}
