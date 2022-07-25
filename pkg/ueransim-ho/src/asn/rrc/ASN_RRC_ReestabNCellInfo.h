/*
 * Generated by asn1c-0.9.29 (http://lionet.info/asn1c)
 * From ASN.1 module "NR-InterNodeDefinitions"
 * 	found in "asn/nr-rrc-15.6.0.asn1"
 * 	`asn1c -fcompound-names -pdu=all -findirect-choice -fno-include-deps -gen-PER -no-gen-OER -no-gen-example -D rrc`
 */

#ifndef	_ASN_RRC_ReestabNCellInfo_H_
#define	_ASN_RRC_ReestabNCellInfo_H_


#include <asn_application.h>

/* Including external dependencies */
#include "ASN_RRC_CellIdentity.h"
#include <BIT_STRING.h>
#include "ASN_RRC_ShortMAC-I.h"
#include <constr_SEQUENCE.h>

#ifdef __cplusplus
extern "C" {
#endif

/* ASN_RRC_ReestabNCellInfo */
typedef struct ASN_RRC_ReestabNCellInfo {
	ASN_RRC_CellIdentity_t	 cellIdentity;
	BIT_STRING_t	 key_gNodeB_Star;
	ASN_RRC_ShortMAC_I_t	 shortMAC_I;
	
	/* Context for parsing across buffer boundaries */
	asn_struct_ctx_t _asn_ctx;
} ASN_RRC_ReestabNCellInfo_t;

/* Implementation */
extern asn_TYPE_descriptor_t asn_DEF_ASN_RRC_ReestabNCellInfo;
extern asn_SEQUENCE_specifics_t asn_SPC_ASN_RRC_ReestabNCellInfo_specs_1;
extern asn_TYPE_member_t asn_MBR_ASN_RRC_ReestabNCellInfo_1[3];

#ifdef __cplusplus
}
#endif

#endif	/* _ASN_RRC_ReestabNCellInfo_H_ */
#include <asn_internal.h>
