//
// This file is a part of UERANSIM open source project.
// Copyright (c) 2021 ALİ GÜNGÖR.
//
// The software and all associated files are licensed under GPL-3.0
// and subject to the terms and conditions defined in LICENSE file.
//

#include "task.hpp"
#include "utils.hpp"

#include <algorithm>

#include <gnb/app/task.hpp>
#include <gnb/rrc/task.hpp>
#include <gnb/sctp/task.hpp>
//kasia
#include <gnb/types.hpp>
#include <iostream>
//

#include <asn/ngap/ASN_NGAP_AMFConfigurationUpdate.h>
#include <asn/ngap/ASN_NGAP_AMFConfigurationUpdateFailure.h>
#include <asn/ngap/ASN_NGAP_AMFName.h>
#include <asn/ngap/ASN_NGAP_BroadcastPLMNItem.h>
#include <asn/ngap/ASN_NGAP_ErrorIndication.h>
#include <asn/ngap/ASN_NGAP_GlobalGNB-ID.h>
#include <asn/ngap/ASN_NGAP_InitiatingMessage.h>
#include <asn/ngap/ASN_NGAP_NGAP-PDU.h>
#include <asn/ngap/ASN_NGAP_NGSetupRequest.h>
#include <asn/ngap/ASN_NGAP_OverloadStartNSSAIItem.h>
#include <asn/ngap/ASN_NGAP_PLMNSupportItem.h>
#include <asn/ngap/ASN_NGAP_ProtocolIE-Field.h>
#include <asn/ngap/ASN_NGAP_ServedGUAMIItem.h>
#include <asn/ngap/ASN_NGAP_SliceSupportItem.h>
#include <asn/ngap/ASN_NGAP_SupportedTAItem.h>
//kasia
#include <asn/ngap/ASN_NGAP_HandoverRequired.h>
#include <asn/ngap/ASN_NGAP_PDUSessionResourceItemHORqd.h>
#include <asn/ngap/ASN_NGAP_SourceNGRANNode-ToTargetNGRANNode-TransparentContainer.h>
#include <utils/octet_string.cpp>
#include "encode.hpp"
#include <asn/rrc/ASN_RRC_HandoverPreparationInformation-IEs.h>
#include <asn/rrc/ASN_RRC_HandoverPreparationInformation.h>
#include <asn/rrc/ASN_RRC_UE-CapabilityRAT-Container.h>
#include <asn/ngap/ASN_NGAP_PDUSessionResourceAdmittedItem.h>
#include <asn/rrc/ASN_RRC_DL-DCCH-MessageType.h>
#include <asn/rrc/ASN_RRC_DL-DCCH-Message.h>
#include <asn/rrc/ASN_RRC_RRCReconfiguration.h>
#include "transport.cpp"

namespace nr::gnb
{

template <typename T>
static void AssignDefaultAmfConfigs(NgapAmfContext *amf, T *msg)
{
    auto *ie = asn::ngap::GetProtocolIe(msg, ASN_NGAP_ProtocolIE_ID_id_AMFName);
    if (ie)
        amf->amfName = asn::GetPrintableString(ie->AMFName);

    ie = asn::ngap::GetProtocolIe(msg, ASN_NGAP_ProtocolIE_ID_id_RelativeAMFCapacity);
    if (ie)
        amf->relativeCapacity = ie->RelativeAMFCapacity;

    ie = asn::ngap::GetProtocolIe(msg, ASN_NGAP_ProtocolIE_ID_id_ServedGUAMIList);
    if (ie)
    {
        utils::ClearAndDelete(amf->servedGuamiList);

        asn::ForeachItem(ie->ServedGUAMIList, [amf](ASN_NGAP_ServedGUAMIItem &item) {
            auto servedGuami = new ServedGuami();
            if (item.backupAMFName)
                servedGuami->backupAmfName = asn::GetPrintableString(*item.backupAMFName);
            ngap_utils::GuamiFromAsn_Ref(item.gUAMI, servedGuami->guami);
            amf->servedGuamiList.push_back(servedGuami);
        });
    }

    ie = asn::ngap::GetProtocolIe(msg, ASN_NGAP_ProtocolIE_ID_id_PLMNSupportList);
    if (ie)
    {
        utils::ClearAndDelete(amf->plmnSupportList);

        asn::ForeachItem(ie->PLMNSupportList, [amf](ASN_NGAP_PLMNSupportItem &item) {
            auto plmnSupport = new PlmnSupport();
            ngap_utils::PlmnFromAsn_Ref(item.pLMNIdentity, plmnSupport->plmn);
            asn::ForeachItem(item.sliceSupportList, [plmnSupport](ASN_NGAP_SliceSupportItem &ssItem) {
                plmnSupport->sliceSupportList.slices.push_back(ngap_utils::SliceSupportFromAsn(ssItem));
            });
            amf->plmnSupportList.push_back(plmnSupport);
        });
    }
}

void NgapTask::handleAssociationSetup(int amfId, int ascId, int inCount, int outCount)
{
    auto *amf = findAmfContext(amfId);
    if (amf != nullptr)
    {
        amf->association.associationId = amfId;
        amf->association.inStreams = inCount;
        amf->association.outStreams = outCount;

        sendNgSetupRequest(amf->ctxId);
    }
}

void NgapTask::handleAssociationShutdown(int amfId)
{
    auto *amf = findAmfContext(amfId);
    if (amf == nullptr)
        return;

    m_logger->err("Association terminated for AMF[%d]", amfId);
    m_logger->debug("Removing AMF context[%d]", amfId);

    amf->state = EAmfState::NOT_CONNECTED;

    auto w = std::make_unique<NmGnbSctp>(NmGnbSctp::CONNECTION_CLOSE);
    w->clientId = amfId;
    m_base->sctpTask->push(std::move(w));

    deleteAmfContext(amfId);
}

void NgapTask::sendNgSetupRequest(int amfId)
{
    m_logger->debug("Sending NG Setup Request");

    auto *amf = findAmfContext(amfId);
    if (amf == nullptr)
        return;

    amf->state = EAmfState::WAITING_NG_SETUP;

    // TODO: this procedure also re-initialises the NGAP UE-related contexts (if any)
    //  and erases all related signalling connections in the two nodes like an NG Reset procedure would do.
    //  More on 38.413 8.7.1.1

    auto *globalGnbId = asn::New<ASN_NGAP_GlobalGNB_ID>();
    globalGnbId->gNB_ID.present = ASN_NGAP_GNB_ID_PR_gNB_ID;
    asn::SetBitString(globalGnbId->gNB_ID.choice.gNB_ID, octet4{m_base->config->getGnbId()},
                      static_cast<size_t>(m_base->config->gnbIdLength));
    asn::SetOctetString3(globalGnbId->pLMNIdentity, ngap_utils::PlmnToOctet3(m_base->config->plmn));

    auto *ieGlobalGnbId = asn::New<ASN_NGAP_NGSetupRequestIEs>();
    ieGlobalGnbId->id = ASN_NGAP_ProtocolIE_ID_id_GlobalRANNodeID;
    ieGlobalGnbId->criticality = ASN_NGAP_Criticality_reject;
    ieGlobalGnbId->value.present = ASN_NGAP_NGSetupRequestIEs__value_PR_GlobalRANNodeID;
    ieGlobalGnbId->value.choice.GlobalRANNodeID.present = ASN_NGAP_GlobalRANNodeID_PR_globalGNB_ID;
    ieGlobalGnbId->value.choice.GlobalRANNodeID.choice.globalGNB_ID = globalGnbId;

    auto *ieRanNodeName = asn::New<ASN_NGAP_NGSetupRequestIEs>();
    ieRanNodeName->id = ASN_NGAP_ProtocolIE_ID_id_RANNodeName;
    ieRanNodeName->criticality = ASN_NGAP_Criticality_ignore;
    ieRanNodeName->value.present = ASN_NGAP_NGSetupRequestIEs__value_PR_RANNodeName;
    asn::SetPrintableString(ieRanNodeName->value.choice.RANNodeName, m_base->config->name);

    auto *broadcastPlmn = asn::New<ASN_NGAP_BroadcastPLMNItem>();
    asn::SetOctetString3(broadcastPlmn->pLMNIdentity, ngap_utils::PlmnToOctet3(m_base->config->plmn));
    for (auto &nssai : m_base->config->nssai.slices)
    {
        auto *item = asn::New<ASN_NGAP_SliceSupportItem>();
        asn::SetOctetString1(item->s_NSSAI.sST, static_cast<uint8_t>(nssai.sst));
        if (nssai.sd.has_value())
        {
            item->s_NSSAI.sD = asn::New<ASN_NGAP_SD_t>();
            asn::SetOctetString3(*item->s_NSSAI.sD, octet3{nssai.sd.value()});
        }
        asn::SequenceAdd(broadcastPlmn->tAISliceSupportList, item);
    }

    auto *supportedTa = asn::New<ASN_NGAP_SupportedTAItem>();
    asn::SetOctetString3(supportedTa->tAC, octet3{m_base->config->tac});
    asn::SequenceAdd(supportedTa->broadcastPLMNList, broadcastPlmn);

    auto *ieSupportedTaList = asn::New<ASN_NGAP_NGSetupRequestIEs>();
    ieSupportedTaList->id = ASN_NGAP_ProtocolIE_ID_id_SupportedTAList;
    ieSupportedTaList->criticality = ASN_NGAP_Criticality_reject;
    ieSupportedTaList->value.present = ASN_NGAP_NGSetupRequestIEs__value_PR_SupportedTAList;
    asn::SequenceAdd(ieSupportedTaList->value.choice.SupportedTAList, supportedTa);

    auto *iePagingDrx = asn::New<ASN_NGAP_NGSetupRequestIEs>();
    iePagingDrx->id = ASN_NGAP_ProtocolIE_ID_id_DefaultPagingDRX;
    iePagingDrx->criticality = ASN_NGAP_Criticality_ignore;
    iePagingDrx->value.present = ASN_NGAP_NGSetupRequestIEs__value_PR_PagingDRX;
    iePagingDrx->value.choice.PagingDRX = ngap_utils::PagingDrxToAsn(m_base->config->pagingDrx);

    auto *pdu = asn::ngap::NewMessagePdu<ASN_NGAP_NGSetupRequest>(
        {ieGlobalGnbId, ieRanNodeName, ieSupportedTaList, iePagingDrx});

    sendNgapNonUe(amfId, pdu);
}

void NgapTask::receiveNgSetupResponse(int amfId, ASN_NGAP_NGSetupResponse *msg)
{
    m_logger->debug("NG Setup Response received");

    auto *amf = findAmfContext(amfId);
    if (amf == nullptr)
        return;

    AssignDefaultAmfConfigs(amf, msg);

    amf->state = EAmfState::CONNECTED;
    m_logger->info("NG Setup procedure is successful");

    if (!m_isInitialized && std::all_of(m_amfCtx.begin(), m_amfCtx.end(),
                                        [](auto &amfCtx) { return amfCtx.second->state == EAmfState::CONNECTED; }))
    {
        m_isInitialized = true;

        auto update = std::make_unique<NmGnbStatusUpdate>(NmGnbStatusUpdate::NGAP_IS_UP);
        update->isNgapUp = true;
        m_base->appTask->push(std::move(update));

        m_base->rrcTask->push(std::make_unique<NmGnbNgapToRrc>(NmGnbNgapToRrc::RADIO_POWER_ON));
    }
}

void NgapTask::receiveNgSetupFailure(int amfId, ASN_NGAP_NGSetupFailure *msg)
{
    auto *amf = findAmfContext(amfId);
    if (amf == nullptr)
        return;

    amf->state = EAmfState::WAITING_NG_SETUP;

    auto *ie = asn::ngap::GetProtocolIe(msg, ASN_NGAP_ProtocolIE_ID_id_Cause);
    if (ie)
        m_logger->err("NG Setup procedure is failed. Cause: %s", ngap_utils::CauseToString(ie->Cause).c_str());
    else
        m_logger->err("NG Setup procedure is failed.");
}

void NgapTask::receiveErrorIndication(int amfId, ASN_NGAP_ErrorIndication *msg)
{
    auto *amf = findAmfContext(amfId);
    if (amf == nullptr)
    {
        m_logger->err("Error indication received with not found AMF context");
        return;
    }

    auto *ie = asn::ngap::GetProtocolIe(msg, ASN_NGAP_ProtocolIE_ID_id_Cause);
    if (ie)
        m_logger->err("Error indication received. Cause: %s", ngap_utils::CauseToString(ie->Cause).c_str());
    else
        m_logger->err("Error indication received.");
}
//kasia
void NgapTask::handleN2Handover(int ueId, int amf_ue_ngap_id, long ran_ue_ngap_id, int targetCellId) //pmq
{
    //m_logger->debug("Sending NG Handover Required");
    sendNgHandoverRequired(ueId, amf_ue_ngap_id, ran_ue_ngap_id, nr::gnb::HandoverType::Intra5GS, 
                                nr::gnb::NgapCause::RadioNetwork_handover_desirable_for_radio_reason, 2, targetCellId);

}
void NgapTask::sendNgHandoverRequired(int ueId, int amf_ue_ngap_id, long ran_ue_ngap_id, HandoverType handover_type, 
                                NgapCause cause, int global_ran_node_id, int targetCellId) // pmq
{
    m_logger->debug("Sending NG Handover Required for UE[%d], AMF_UE_NGAP_ID[%d], RAN_UE_NGAP_ID[%d] and TARGET_CELL_ID[%d]", ueId, amf_ue_ngap_id, ran_ue_ngap_id, targetCellId);
    

    auto *ieHandoverType = asn::New<ASN_NGAP_HandoverRequiredIEs>();
    ieHandoverType->id = ASN_NGAP_ProtocolIE_ID_id_HandoverType;
    ieHandoverType->criticality = ASN_NGAP_Criticality_reject;
    ieHandoverType->value.present = ASN_NGAP_HandoverRequiredIEs__value_PR_HandoverType;
    ieHandoverType->value.choice.HandoverType = ASN_NGAP_HandoverType_intra5gs; //hardoced, the only considered type

    auto *ieCause = asn::New<ASN_NGAP_HandoverRequiredIEs>();
    ieCause->id = ASN_NGAP_ProtocolIE_ID_id_Cause; 
    ieCause->criticality = ASN_NGAP_Criticality_ignore;
    ieCause->value.present = ASN_NGAP_HandoverRequiredIEs__value_PR_Cause;
    ngap_utils::ToCauseAsn_Ref(nr::gnb::NgapCause::RadioNetwork_handover_desirable_for_radio_reason, ieCause->value.choice.Cause); //hardcoded


    auto *globalGnbId = asn::New<ASN_NGAP_GlobalGNB_ID>();
    asn::SetOctetString3(globalGnbId->pLMNIdentity, ngap_utils::PlmnToOctet3(m_base->config->plmn));
    globalGnbId->gNB_ID.present = ASN_NGAP_GNB_ID_PR_gNB_ID;
    asn::SetBitString(globalGnbId->gNB_ID.choice.gNB_ID, octet4{targetCellId}, // Its not always ok. octet4 is no directly mapping with int values
                      static_cast<size_t>(m_base->config->gnbIdLength));

    auto *targetRANNodeID = asn::New<ASN_NGAP_TargetRANNodeID>();
    targetRANNodeID->globalRANNodeID.present = ASN_NGAP_GlobalRANNodeID_PR_globalGNB_ID;
    targetRANNodeID->globalRANNodeID.choice.globalGNB_ID = globalGnbId;
    asn::SetOctetString3(targetRANNodeID->selectedTAI.pLMNIdentity, ngap_utils::PlmnToOctet3(m_base->config->plmn));
    asn::SetOctetString3(targetRANNodeID->selectedTAI.tAC, octet3{m_base->config->tac});

    auto *ieTarget = asn::New<ASN_NGAP_HandoverRequiredIEs>();
    ieTarget->id = ASN_NGAP_ProtocolIE_ID_id_TargetID;
    ieTarget->criticality = ASN_NGAP_Criticality_reject;
    ieTarget->value.present = ASN_NGAP_HandoverRequiredIEs__value_PR_TargetID;
    ieTarget->value.choice.TargetID.present = ASN_NGAP_TargetID_PR_targetRANNodeID;
    ieTarget->value.choice.TargetID.choice.targetRANNodeID = targetRANNodeID;

    auto *PDUSessionResourceItem = asn::New<ASN_NGAP_PDUSessionResourceItemHORqd>();
    PDUSessionResourceItem->pDUSessionID = findPduByUeId(ueId);

    auto transfer = OctetString::FromHex("40"); // hardcoded, means XnAP is not available
    asn::SetOctetString(PDUSessionResourceItem->handoverRequiredTransfer, transfer); 

    auto *iePDUSessionResourceList = asn::New<ASN_NGAP_HandoverRequiredIEs>();
    iePDUSessionResourceList->id = ASN_NGAP_ProtocolIE_ID_id_PDUSessionResourceListHORqd;
    iePDUSessionResourceList->criticality = ASN_NGAP_Criticality_reject;
    iePDUSessionResourceList->value.present = ASN_NGAP_HandoverRequiredIEs__value_PR_PDUSessionResourceListHORqd;

    asn::SequenceAdd(iePDUSessionResourceList->value.choice.PDUSessionResourceListHORqd, PDUSessionResourceItem);
    auto *sourceToTargetContainer = asn::New<ASN_NGAP_HandoverRequiredIEs>();
    sourceToTargetContainer->id = ASN_NGAP_ProtocolIE_ID_id_SourceToTarget_TransparentContainer;
    sourceToTargetContainer->criticality = ASN_NGAP_Criticality_reject;
    sourceToTargetContainer->value.present = ASN_NGAP_HandoverRequiredIEs__value_PR_SourceToTarget_TransparentContainer;

    auto hex = OctetString::FromHex("400300001100000100010002F8390000432100000002F839000000001000000A");
    asn::SetOctetString(sourceToTargetContainer->value.choice.SourceToTarget_TransparentContainer, hex);
    auto *pdu = asn::ngap::NewMessagePdu<ASN_NGAP_HandoverRequired>(
        {/*ieAmfUeId, ieRanUeId, */ieHandoverType, ieCause, ieTarget, iePDUSessionResourceList, sourceToTargetContainer});

    sendNgapUeAssociated(ueId, pdu);


}

void NgapTask::receiveHandoverRequest(int amfId, ASN_NGAP_HandoverRequest *msg)
{
    m_logger->debug("Handover Request Recived");

    auto *amf = findAmfContext(amfId);
    if (amf == nullptr)
        return;

    sendHandoverRequestAck(amfId, msg);
}

void NgapTask::sendHandoverRequestAck(int amfId, ASN_NGAP_HandoverRequest *msg)
{
    m_logger->debug("Sending Handover Request Acknowledge");

    auto *amf = findAmfContext(amfId);
    if (amf == nullptr)
        return;

    auto *ieAmfUeId = asn::New<ASN_NGAP_HandoverRequestAcknowledgeIEs>();
    ieAmfUeId->id = ASN_NGAP_ProtocolIE_ID_id_AMF_UE_NGAP_ID;
    ieAmfUeId->criticality = ASN_NGAP_Criticality_ignore;
    ieAmfUeId->value.present = ASN_NGAP_HandoverRequestAcknowledgeIEs__value_PR_AMF_UE_NGAP_ID;
    auto *ie = asn::ngap::GetProtocolIe(msg, ASN_NGAP_ProtocolIE_ID_id_AMF_UE_NGAP_ID);
    ieAmfUeId->value.choice.AMF_UE_NGAP_ID = ie->AMF_UE_NGAP_ID;

    auto *ieRanUeId = asn::New<ASN_NGAP_HandoverRequestAcknowledgeIEs>();
    ieRanUeId->id = ASN_NGAP_ProtocolIE_ID_id_RAN_UE_NGAP_ID;
    ieRanUeId->criticality = ASN_NGAP_Criticality_ignore;
    ieRanUeId->value.present = ASN_NGAP_HandoverRequestAcknowledgeIEs__value_PR_RAN_UE_NGAP_ID;
    // auto *ieRanUeNgapId = asn::ngap::GetProtocolIe(msg, ASN_NGAP_ProtocolIE_ID_id_RAN_UE_NGAP_ID);
    ieRanUeId->value.choice.RAN_UE_NGAP_ID = 1; //hardcoded, considers only UE with RAN_UE_NGAP_ID=1

    auto *iePduAdmitted = asn::New<ASN_NGAP_HandoverRequestAcknowledgeIEs>();
    iePduAdmitted->id = ASN_NGAP_ProtocolIE_ID_id_PDUSessionResourceAdmittedList;
    iePduAdmitted->criticality = ASN_NGAP_Criticality_ignore;
    iePduAdmitted->value.present = ASN_NGAP_HandoverRequestAcknowledgeIEs__value_PR_PDUSessionResourceAdmittedList;
    auto *iePduAdmittedItem = asn::New<ASN_NGAP_PDUSessionResourceAdmittedItem>();
    iePduAdmittedItem->pDUSessionID = 1; //hardoded, considers only pdu session with id 1
    auto transfer = OctetString::FromHex("0007C0AC100301000000010001"); //hardcoded
    asn::SetOctetString(iePduAdmittedItem->handoverRequestAcknowledgeTransfer, transfer);
    asn::SequenceAdd(iePduAdmitted->value.choice.PDUSessionResourceAdmittedList, iePduAdmittedItem);


    auto *ieTargetContainer = asn::New<ASN_NGAP_HandoverRequestAcknowledgeIEs>();
    ieTargetContainer->id = ASN_NGAP_ProtocolIE_ID_id_TargetToSource_TransparentContainer;
    ieTargetContainer->criticality = ASN_NGAP_Criticality_reject;
    ieTargetContainer->value.present = ASN_NGAP_HandoverRequestAcknowledgeIEs__value_PR_TargetToSource_TransparentContainer;
    auto container = OctetString::FromHex("00010000"); // hardcoded
    asn::SetOctetString(ieTargetContainer->value.choice.TargetToSource_TransparentContainer, container);

    auto *pdu_non = asn::ngap::NewMessagePdu<ASN_NGAP_HandoverRequestAcknowledge>(
        {ieAmfUeId, ieRanUeId, iePduAdmitted, ieTargetContainer});

    sendNgapNonUeResponse(amfId, pdu_non);
    
    sendHandoverNotify(amfId, msg);
    
}

void NgapTask::receiveHandoverCommand(int amfId, ASN_NGAP_HandoverCommand *msg){
    m_logger->debug("Handover Command Recived");

    auto *ieAmfUeNgapId = asn::ngap::GetProtocolIe(msg, ASN_NGAP_ProtocolIE_ID_id_AMF_UE_NGAP_ID);
    int64_t amfUeNgapId;
    amfUeNgapId = asn::GetSigned64(ieAmfUeNgapId->AMF_UE_NGAP_ID);
    //std::cout<<amfUeNgapId<<std::endl;

    auto *ue = findUeByAmfId(amfUeNgapId);
    if (ue == nullptr)
        return;
    auto w1 = std::make_unique<NmGnbNgapToRrc>(NmGnbNgapToRrc::RRC_RECONFIGURATION);
    w1->ueId = ue->ctxId;
    m_base->rrcTask->push(std::move(w1));

     
}

void NgapTask::sendHandoverNotify(int amfId, ASN_NGAP_HandoverRequest *msg){
    m_logger->debug("Handover notify");
    
    //ASN_NGAP_HandoverNotifyIEs
    auto *ieAmfUeId = asn::New<ASN_NGAP_HandoverNotifyIEs>();
    ieAmfUeId->id = ASN_NGAP_ProtocolIE_ID_id_AMF_UE_NGAP_ID;
    ieAmfUeId->criticality = ASN_NGAP_Criticality_reject;
    ieAmfUeId->value.present = ASN_NGAP_HandoverNotifyIEs__value_PR_AMF_UE_NGAP_ID;
    auto *ie = asn::ngap::GetProtocolIe(msg, ASN_NGAP_ProtocolIE_ID_id_AMF_UE_NGAP_ID);
    ieAmfUeId->value.choice.AMF_UE_NGAP_ID = ie->AMF_UE_NGAP_ID;

    auto *ieRanUeId = asn::New<ASN_NGAP_HandoverNotifyIEs>();
    ieRanUeId->id = ASN_NGAP_ProtocolIE_ID_id_RAN_UE_NGAP_ID;
    ieRanUeId->criticality = ASN_NGAP_Criticality_reject;
    ieRanUeId->value.present = ASN_NGAP_HandoverNotifyIEs__value_PR_RAN_UE_NGAP_ID;
    //auto *ie2 = asn::ngap::GetProtocolIe(msg, ASN_NGAP_ProtocolIE_ID_id_RAN_UE_NGAP_ID);
    ieRanUeId->value.choice.RAN_UE_NGAP_ID = 1;

    
    auto *pdu = asn::ngap::NewMessagePdu<ASN_NGAP_HandoverNotify>(
        {ieAmfUeId,ieRanUeId});

     asn::ngap::AddProtocolIeIfUsable(
            *pdu, asn_DEF_ASN_NGAP_UserLocationInformation, ASN_NGAP_ProtocolIE_ID_id_UserLocationInformation,
            FindCriticalityOfUserIe(pdu, ASN_NGAP_ProtocolIE_ID_id_UserLocationInformation), [this](void *mem) {
                auto *loc = reinterpret_cast<ASN_NGAP_UserLocationInformation *>(mem);
                loc->present = ASN_NGAP_UserLocationInformation_PR_userLocationInformationNR;
                loc->choice.userLocationInformationNR = asn::New<ASN_NGAP_UserLocationInformationNR>();

                auto &nr = loc->choice.userLocationInformationNR;
                nr->timeStamp = asn::New<ASN_NGAP_TimeStamp_t>();

                ngap_utils::ToPlmnAsn_Ref(m_base->config->plmn, nr->nR_CGI.pLMNIdentity);
                asn::SetBitStringLong<36>(m_base->config->nci, nr->nR_CGI.nRCellIdentity);
                ngap_utils::ToPlmnAsn_Ref(m_base->config->plmn, nr->tAI.pLMNIdentity);
                asn::SetOctetString3(nr->tAI.tAC, octet3{m_base->config->tac});
                asn::SetOctetString4(*nr->timeStamp, octet4{utils::CurrentTimeStamp().seconds32()});
            });

    sendNgapNonUeResponse(amfId, pdu);

}



void NgapTask::sendErrorIndication(int amfId, NgapCause cause, int ueId)
{
    auto ieCause = asn::New<ASN_NGAP_ErrorIndicationIEs>();
    ieCause->id = ASN_NGAP_ProtocolIE_ID_id_Cause;
    ieCause->criticality = ASN_NGAP_Criticality_ignore;
    ieCause->value.present = ASN_NGAP_ErrorIndicationIEs__value_PR_Cause;
    ngap_utils::ToCauseAsn_Ref(cause, ieCause->value.choice.Cause);

    m_logger->warn("Sending an error indication with cause: %s",
                   ngap_utils::CauseToString(ieCause->value.choice.Cause).c_str());

    auto *pdu = asn::ngap::NewMessagePdu<ASN_NGAP_ErrorIndication>({ieCause});

    if (ueId > 0)
        sendNgapUeAssociated(ueId, pdu);
    else
        sendNgapNonUe(amfId, pdu);
}

void NgapTask::receiveAmfConfigurationUpdate(int amfId, ASN_NGAP_AMFConfigurationUpdate *msg)
{
    m_logger->debug("AMF configuration update received");

    auto *amf = findAmfContext(amfId);
    if (amf == nullptr)
        return;

    bool tnlModified = false;

    auto *ie = asn::ngap::GetProtocolIe(msg, ASN_NGAP_ProtocolIE_ID_id_AMF_TNLAssociationToAddList);
    if (ie && ie->AMF_TNLAssociationToAddList.list.count > 0)
        tnlModified = true;

    ie = asn::ngap::GetProtocolIe(msg, ASN_NGAP_ProtocolIE_ID_id_AMF_TNLAssociationToRemoveList);
    if (ie && ie->AMF_TNLAssociationToRemoveList.list.count > 0)
        tnlModified = true;

    ie = asn::ngap::GetProtocolIe(msg, ASN_NGAP_ProtocolIE_ID_id_AMF_TNLAssociationToUpdateList);
    if (ie && ie->AMF_TNLAssociationToUpdateList.list.count > 0)
        tnlModified = true;

    // TODO: AMF TNL modification is not supported
    if (tnlModified)
    {
        m_logger->err("TNL modification is not supported, rejecting AMF configuration update");

        auto *ieCause = asn::New<ASN_NGAP_AMFConfigurationUpdateFailureIEs>();
        ieCause->id = ASN_NGAP_ProtocolIE_ID_id_Cause;
        ieCause->criticality = ASN_NGAP_Criticality_ignore;
        ieCause->value.present = ASN_NGAP_AMFConfigurationUpdateFailureIEs__value_PR_Cause;
        ngap_utils::ToCauseAsn_Ref(NgapCause::Transport_unspecified, ieCause->value.choice.Cause);

        auto *pdu = asn::ngap::NewMessagePdu<ASN_NGAP_AMFConfigurationUpdateFailure>({ieCause});
        sendNgapNonUe(amfId, pdu);
    }
    else
    {
        AssignDefaultAmfConfigs(amf, msg);

        auto *ieList = asn::New<ASN_NGAP_AMFConfigurationUpdateAcknowledgeIEs>();
        ieList->id = ASN_NGAP_ProtocolIE_ID_id_AMF_TNLAssociationSetupList;
        ieList->criticality = ASN_NGAP_Criticality_ignore;
        ieList->value.present = ASN_NGAP_AMFConfigurationUpdateAcknowledgeIEs__value_PR_AMF_TNLAssociationSetupList;

        auto *pdu = asn::ngap::NewMessagePdu<ASN_NGAP_AMFConfigurationUpdateAcknowledge>({ieList});
        sendNgapNonUe(amfId, pdu);
    }
}

void NgapTask::receiveOverloadStart(int amfId, ASN_NGAP_OverloadStart *msg)
{
    m_logger->debug("AMF overload start received");

    auto *amf = findAmfContext(amfId);
    if (amf == nullptr)
        return;

    amf->overloadInfo = {};
    amf->overloadInfo.status = EOverloadStatus::OVERLOADED;

    auto *ie = asn::ngap::GetProtocolIe(msg, ASN_NGAP_ProtocolIE_ID_id_AMFOverloadResponse);
    if (ie && ie->OverloadResponse.present == ASN_NGAP_OverloadResponse_PR_overloadAction)
    {
        switch (ie->OverloadResponse.choice.overloadAction)
        {
        case ASN_NGAP_OverloadAction_reject_non_emergency_mo_dt:
            amf->overloadInfo.indication.action = EOverloadAction::REJECT_NON_EMERGENCY_MO_DATA;
            break;
        case ASN_NGAP_OverloadAction_reject_rrc_cr_signalling:
            amf->overloadInfo.indication.action = EOverloadAction::REJECT_SIGNALLING;
            break;
        case ASN_NGAP_OverloadAction_permit_emergency_sessions_and_mobile_terminated_services_only:
            amf->overloadInfo.indication.action = EOverloadAction::ONLY_EMERGENCY_AND_MT;
            break;
        case ASN_NGAP_OverloadAction_permit_high_priority_sessions_and_mobile_terminated_services_only:
            amf->overloadInfo.indication.action = EOverloadAction::ONLY_HIGH_PRI_AND_MT;
            break;
        default:
            m_logger->warn("AMF overload action [%d] could not understand",
                           (int)ie->OverloadResponse.choice.overloadAction);
            break;
        }
    }

    ie = asn::ngap::GetProtocolIe(msg, ASN_NGAP_ProtocolIE_ID_id_AMFTrafficLoadReductionIndication);
    if (ie)
        amf->overloadInfo.indication.loadReductionPerc = static_cast<int>(ie->TrafficLoadReductionIndication);

    ie = asn::ngap::GetProtocolIe(msg, ASN_NGAP_ProtocolIE_ID_id_OverloadStartNSSAIList);
    if (ie)
    {
        // TODO
        /*asn::ForeachItem(ie->OverloadStartNSSAIList, [](auto &item) {
            item.sliceOverloadList;
        });*/
    }
}

void NgapTask::receiveOverloadStop(int amfId, ASN_NGAP_OverloadStop *msg)
{
    m_logger->debug("AMF overload stop received");

    // TODO
}

} // namespace nr::gnb
