package ts29338

import dia "github.com/fkgi/diameter"

const (
	// DiameterErrorUserUnknown is Result-Code 5001
	DiameterErrorUserUnknown uint32 = dia.ResultOffset + 5001
	// DiameterErrorAbsentUser is Result-Code 5550
	DiameterErrorAbsentUser uint32 = dia.ResultOffset + 5550
	// DiameterErrorUserBusyForMtSms is Result-Code 5551
	DiameterErrorUserBusyForMtSms uint32 = dia.ResultOffset + 5551
	// DiameterErrorFacilityNotSupported is Result-Code 5552
	DiameterErrorFacilityNotSupported uint32 = dia.ResultOffset + 5552
	// DiameterErrorIlleagalUser is Result-Code 5553
	DiameterErrorIlleagalUser uint32 = dia.ResultOffset + 5553
	// DiameterErrorIlleagalEquipment is Result-Code 5554
	DiameterErrorIlleagalEquipment uint32 = dia.ResultOffset + 5554
	// DiameterErrorSmDeliveryFailure is Result-Code 5555
	DiameterErrorSmDeliveryFailure uint32 = dia.ResultOffset + 5555
	// DiameterErrorServiceNotSubscribed is Result-Code 5556
	DiameterErrorServiceNotSubscribed uint32 = dia.ResultOffset + 5556
	// DiameterErrorServiceBarred is Result-Code 5557
	DiameterErrorServiceBarred uint32 = dia.ResultOffset + 5557
	// DiameterErrorMwdListFull is Result-Code 5558
	DiameterErrorMwdListFull uint32 = dia.ResultOffset + 5558
)
