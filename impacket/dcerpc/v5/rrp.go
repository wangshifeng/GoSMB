// SECUREAUTH LABS. Copyright 2018 SecureAuth Corporation. All rights reserved.
//
// This software is provided under under a slightly modified version
// of the Apache Software License. See the accompanying LICENSE file
// for more information.
//
// Author: Alberto Solino (@agsolino)
//
// Description:
//   [MS-RRP] Interface implementation
//
//   Best way to learn how to use these calls is to grab the protocol standard
//   so you understand what the call does, and then read the test case located
//   at https://github.com/SecureAuthCorp/impacket/tree/master/tests/SMB_RPC
//
//   Some calls have helper functions, which makes it even easier to use.
//   They are located at the end of this file. 
//   Helper functions start with "h"<name of the call>.
//   There are test cases for them too. 
//


from struct import unpack, pack

from impacket.dcerpc.v5.ndr import NDRCALL, NDRSTRUCT, NDRPOINTER, NDRUniConformantVaryingArray, NDRUniConformantArray
from impacket.dcerpc.v5.dtypes import DWORD, UUID, ULONG, LPULONG, BOOLEAN, SECURITY_INFORMATION, PFILETIME, \
    RPC_UNICODE_STRING, FILETIME, NULL, MAXIMUM_ALLOWED, OWNER_SECURITY_INFORMATION, PWCHAR, PRPC_UNICODE_STRING
from impacket.dcerpc.v5.rpcrt import DCERPCException
from impacket import system_errors, LOG
from impacket.uuid import uuidtup_to_bin

MSRPC_UUID_RRP = uuidtup_to_bin(('338CD001-2244-31F1-AAAA-900038001003', '1.0'))

 type DCERPCSessionError struct { // DCERPCException:
     func (self TYPE) __init__(error_string=nil, error_code=nil, packet=nil interface{}){
        DCERPCException.__init__(self, error_string, error_code, packet)

     func __str__( self  interface{}){
        key = self.error_code
        if key in system_errors.ERROR_MESSAGES {
            error_msg_short = system_errors.ERROR_MESSAGES[key][0]
            error_msg_verbose = system_errors.ERROR_MESSAGES[key][1] 
            return 'RRP SessionError: code: 0x%x - %s - %s' % (self.error_code, error_msg_short, error_msg_verbose)
        } else  {
            return 'RRP SessionError: unknown error code: 0x%x' % self.error_code

//###############################################################################
// CONSTANTS
//###############################################################################
// 2.2.2 PREGISTRY_SERVER_NAME
PREGISTRY_SERVER_NAME = PWCHAR

// 2.2.3 error_status_t
error_status_t = ULONG

// 2.2.5 RRP_UNICODE_STRING
RRP_UNICODE_STRING = RPC_UNICODE_STRING
PRRP_UNICODE_STRING = PRPC_UNICODE_STRING

// 2.2.4 REGSAM
REGSAM = ULONG

KEY_QUERY_VALUE        = 0x00000001
KEY_SET_VALUE          = 0x00000002
KEY_CREATE_SUB_KEY     = 0x00000004
KEY_ENUMERATE_SUB_KEYS = 0x00000008
KEY_CREATE_LINK        = 0x00000020
KEY_WOW64_64KEY        = 0x00000100
KEY_WOW64_32KEY        = 0x00000200

REG_BINARY              = 3
REG_DWORD               = 4
REG_DWORD_LITTLE_ENDIAN = 4
REG_DWORD_BIG_ENDIAN    = 5
REG_EXPAND_SZ           = 2
REG_LINK                = 6
REG_MULTI_SZ            = 7
REG_NONE                = 0
REG_QWORD               = 11
REG_QWORD_LITTLE_ENDIAN = 11
REG_SZ                  = 1 

// 3.1.5.7 BaseRegCreateKey (Opnum 6)
REG_CREATED_NEW_KEY     = 0x00000001
REG_OPENED_EXISTING_KEY = 0x00000002

// 3.1.5.19 BaseRegRestoreKey (Opnum 19)
// Flags
REG_WHOLE_HIVE_VOLATILE = 0x00000001
REG_REFRESH_HIVE        = 0x00000002
REG_NO_LAZY_FLUSH       = 0x00000004
REG_FORCE_RESTORE       = 0x00000008

//###############################################################################
// STRUCTURES
//###############################################################################
// 2.2.1 RPC_HKEY
 type RPC_HKEY struct { // NDRSTRUCT:  (
        ('context_handle_attributes',ULONG),
        ('context_handle_uuid',UUID),
    }
     func (self TYPE) __init__(data = nil,isNDR64 = false interface{}){
        NDRSTRUCT.__init__(self, data, isNDR64)
        self.context_handle_uuid = "\x00"*20

// 2.2.6 RVALENT
 type RVALENT struct { // NDRSTRUCT:  (
        ('ve_valuename',PRRP_UNICODE_STRING),
        ('ve_valuelen',DWORD),
        ('ve_valueptr',DWORD),
        ('ve_type',DWORD),
    }

 type RVALENT_ARRAY struct { // NDRUniConformantVaryingArray:
    item = RVALENT

// 2.2.9 RPC_SECURITY_DESCRIPTOR
 type BYTE_ARRAY struct { // NDRUniConformantVaryingArray:
    pass

 type PBYTE_ARRAY struct { // NDRPOINTER:
    referent = (
        ('Data', BYTE_ARRAY),
    }

 type RPC_SECURITY_DESCRIPTOR struct { // NDRSTRUCT:  (
        ('lpSecurityDescriptor',PBYTE_ARRAY),
        ('cbInSecurityDescriptor',DWORD),
        ('cbOutSecurityDescriptor',DWORD),
    }

// 2.2.8 RPC_SECURITY_ATTRIBUTES
 type RPC_SECURITY_ATTRIBUTES struct { // NDRSTRUCT:  (
        ('nLength',DWORD),
        ('RpcSecurityDescriptor',RPC_SECURITY_DESCRIPTOR),
        ('bInheritHandle',BOOLEAN),
    }

 type PRPC_SECURITY_ATTRIBUTES struct { // NDRPOINTER:
    referent = (
        ('Data', RPC_SECURITY_ATTRIBUTES),
    }

//###############################################################################
// RPC CALLS
//###############################################################################
// 3.1.5.1 OpenClassesRoot (Opnum 0)
 type OpenClassesRoot struct { // NDRCALL:
    opnum = 0 (
       ('ServerName', PREGISTRY_SERVER_NAME),
       ('samDesired', REGSAM),
    }

 type OpenClassesRootResponse struct { // NDRCALL: (
       ('phKey', RPC_HKEY),
       ('ErrorCode', error_status_t),
    }

// 3.1.5.2 OpenCurrentUser (Opnum 1)
 type OpenCurrentUser struct { // NDRCALL:
    opnum = 1 (
       ('ServerName', PREGISTRY_SERVER_NAME),
       ('samDesired', REGSAM),
    }

 type OpenCurrentUserResponse struct { // NDRCALL: (
       ('phKey', RPC_HKEY),
       ('ErrorCode', error_status_t),
    }

// 3.1.5.3 OpenLocalMachine (Opnum 2)
 type OpenLocalMachine struct { // NDRCALL:
    opnum = 2 (
       ('ServerName', PREGISTRY_SERVER_NAME),
       ('samDesired', REGSAM),
    }

 type OpenLocalMachineResponse struct { // NDRCALL: (
       ('phKey', RPC_HKEY),
       ('ErrorCode', error_status_t),
    }

// 3.1.5.4 OpenPerformanceData (Opnum 3)
 type OpenPerformanceData struct { // NDRCALL:
    opnum = 3 (
       ('ServerName', PREGISTRY_SERVER_NAME),
       ('samDesired', REGSAM),
    }

 type OpenPerformanceDataResponse struct { // NDRCALL: (
       ('phKey', RPC_HKEY),
       ('ErrorCode', error_status_t),
    }

// 3.1.5.5 OpenUsers (Opnum 4)
 type OpenUsers struct { // NDRCALL:
    opnum = 4 (
       ('ServerName', PREGISTRY_SERVER_NAME),
       ('samDesired', REGSAM),
    }

 type OpenUsersResponse struct { // NDRCALL: (
       ('phKey', RPC_HKEY),
       ('ErrorCode', error_status_t),
    }

// 3.1.5.6 BaseRegCloseKey (Opnum 5)
 type BaseRegCloseKey struct { // NDRCALL:
    opnum = 5 (
       ('hKey', RPC_HKEY),
    }

 type BaseRegCloseKeyResponse struct { // NDRCALL: (
       ('hKey', RPC_HKEY),
       ('ErrorCode', error_status_t),
    }

// 3.1.5.7 BaseRegCreateKey (Opnum 6)
 type BaseRegCreateKey struct { // NDRCALL:
    opnum = 6 (
       ('hKey', RPC_HKEY),
       ('lpSubKey', RRP_UNICODE_STRING),
       ('lpClass', RRP_UNICODE_STRING),
       ('dwOptions', DWORD),
       ('samDesired', REGSAM),
       ('lpSecurityAttributes', PRPC_SECURITY_ATTRIBUTES),
       ('lpdwDisposition', LPULONG),
    }

 type BaseRegCreateKeyResponse struct { // NDRCALL: (
       ('phkResult', RPC_HKEY),
       ('lpdwDisposition', LPULONG),
       ('ErrorCode', error_status_t),
    }

// 3.1.5.8 BaseRegDeleteKey (Opnum 7)
 type BaseRegDeleteKey struct { // NDRCALL:
    opnum = 7 (
       ('hKey', RPC_HKEY),
       ('lpSubKey', RRP_UNICODE_STRING),
    }

 type BaseRegDeleteKeyResponse struct { // NDRCALL: (
       ('ErrorCode', error_status_t),
    }

// 3.1.5.9 BaseRegDeleteValue (Opnum 8)
 type BaseRegDeleteValue struct { // NDRCALL:
    opnum = 8 (
       ('hKey', RPC_HKEY),
       ('lpValueName', RRP_UNICODE_STRING),
    }

 type BaseRegDeleteValueResponse struct { // NDRCALL: (
       ('ErrorCode', error_status_t),
    }

// 3.1.5.10 BaseRegEnumKey (Opnum 9)
 type BaseRegEnumKey struct { // NDRCALL:
    opnum = 9 (
       ('hKey', RPC_HKEY),
       ('dwIndex', DWORD),
       ('lpNameIn', RRP_UNICODE_STRING),
       ('lpClassIn', PRRP_UNICODE_STRING),
       ('lpftLastWriteTime', PFILETIME),
    }

 type BaseRegEnumKeyResponse struct { // NDRCALL: (
       ('lpNameOut', RRP_UNICODE_STRING),
       ('lplpClassOut', PRRP_UNICODE_STRING),
       ('lpftLastWriteTime', PFILETIME),
       ('ErrorCode', error_status_t),
    }

// 3.1.5.11 BaseRegEnumValue (Opnum 10)
 type BaseRegEnumValue struct { // NDRCALL:
    opnum = 10 (
       ('hKey', RPC_HKEY),
       ('dwIndex', DWORD),
       ('lpValueNameIn', RRP_UNICODE_STRING),
       ('lpType', LPULONG),
       ('lpData', PBYTE_ARRAY),
       ('lpcbData', LPULONG),
       ('lpcbLen', LPULONG),
    }

 type BaseRegEnumValueResponse struct { // NDRCALL: (
       ('lpValueNameOut', RRP_UNICODE_STRING),
       ('lpType', LPULONG),
       ('lpData', PBYTE_ARRAY),
       ('lpcbData', LPULONG),
       ('lpcbLen', LPULONG),
       ('ErrorCode', error_status_t),
    }

// 3.1.5.12 BaseRegFlushKey (Opnum 11)
 type BaseRegFlushKey struct { // NDRCALL:
    opnum = 11 (
       ('hKey', RPC_HKEY),
    }

 type BaseRegFlushKeyResponse struct { // NDRCALL: (
       ('ErrorCode', error_status_t),
    }

// 3.1.5.13 BaseRegGetKeySecurity (Opnum 12)
 type BaseRegGetKeySecurity struct { // NDRCALL:
    opnum = 12 (
       ('hKey', RPC_HKEY),
       ('SecurityInformation', SECURITY_INFORMATION),
       ('pRpcSecurityDescriptorIn', RPC_SECURITY_DESCRIPTOR),
    }

 type BaseRegGetKeySecurityResponse struct { // NDRCALL: (
       ('pRpcSecurityDescriptorOut', RPC_SECURITY_DESCRIPTOR),
       ('ErrorCode', error_status_t),
    }

// 3.1.5.14 BaseRegLoadKey (Opnum 13)
 type BaseRegLoadKey struct { // NDRCALL:
    opnum = 13 (
       ('hKey', RPC_HKEY),
       ('lpSubKey', RRP_UNICODE_STRING),
       ('lpFile', RRP_UNICODE_STRING),
    }

 type BaseRegLoadKeyResponse struct { // NDRCALL: (
       ('ErrorCode', error_status_t),
    }

// 3.1.5.15 BaseRegOpenKey (Opnum 15)
 type BaseRegOpenKey struct { // NDRCALL:
    opnum = 15 (
       ('hKey', RPC_HKEY),
       ('lpSubKey', RRP_UNICODE_STRING),
       ('dwOptions', DWORD),
       ('samDesired', REGSAM),
    }

 type BaseRegOpenKeyResponse struct { // NDRCALL: (
       ('phkResult', RPC_HKEY),
       ('ErrorCode', error_status_t),
    }

// 3.1.5.16 BaseRegQueryInfoKey (Opnum 16)
 type BaseRegQueryInfoKey struct { // NDRCALL:
    opnum = 16 (
       ('hKey', RPC_HKEY),
       ('lpClassIn', RRP_UNICODE_STRING),
    }

 type BaseRegQueryInfoKeyResponse struct { // NDRCALL: (
       ('lpClassOut', RPC_UNICODE_STRING),
       ('lpcSubKeys', DWORD),
       ('lpcbMaxSubKeyLen', DWORD),
       ('lpcbMaxClassLen', DWORD),
       ('lpcValues', DWORD),
       ('lpcbMaxValueNameLen', DWORD),
       ('lpcbMaxValueLen', DWORD),
       ('lpcbSecurityDescriptor', DWORD),
       ('lpftLastWriteTime', FILETIME),
       ('ErrorCode', error_status_t),
    }

// 3.1.5.17 BaseRegQueryValue (Opnum 17)
 type BaseRegQueryValue struct { // NDRCALL:
    opnum = 17 (
       ('hKey', RPC_HKEY),
       ('lpValueName', RRP_UNICODE_STRING),
       ('lpType', LPULONG),
       ('lpData', PBYTE_ARRAY),
       ('lpcbData', LPULONG),
       ('lpcbLen', LPULONG),
    }

 type BaseRegQueryValueResponse struct { // NDRCALL: (
       ('lpType', LPULONG),
       ('lpData', PBYTE_ARRAY),
       ('lpcbData', LPULONG),
       ('lpcbLen', LPULONG),
       ('ErrorCode', error_status_t),
    }

// 3.1.5.18 BaseRegReplaceKey (Opnum 18)
 type BaseRegReplaceKey struct { // NDRCALL:
    opnum = 18 (
       ('hKey', RPC_HKEY),
       ('lpSubKey', RRP_UNICODE_STRING),
       ('lpNewFile', RRP_UNICODE_STRING),
       ('lpOldFile', RRP_UNICODE_STRING),
    }

 type BaseRegReplaceKeyResponse struct { // NDRCALL: (
       ('ErrorCode', error_status_t),
    }

// 3.1.5.19 BaseRegRestoreKey (Opnum 19)
 type BaseRegRestoreKey struct { // NDRCALL:
    opnum = 19 (
       ('hKey', RPC_HKEY),
       ('lpFile', RRP_UNICODE_STRING),
       ('Flags', DWORD),
    }

 type BaseRegRestoreKeyResponse struct { // NDRCALL: (
       ('ErrorCode', error_status_t),
    }

// 3.1.5.20 BaseRegSaveKey (Opnum 20)
 type BaseRegSaveKey struct { // NDRCALL:
    opnum = 20 (
       ('hKey', RPC_HKEY),
       ('lpFile', RRP_UNICODE_STRING),
       ('pSecurityAttributes', PRPC_SECURITY_ATTRIBUTES),
    }

 type BaseRegSaveKeyResponse struct { // NDRCALL: (
       ('ErrorCode', error_status_t),
    }

// 3.1.5.21 BaseRegSetKeySecurity (Opnum 21)
 type BaseRegSetKeySecurity struct { // NDRCALL:
    opnum = 21 (
       ('hKey', RPC_HKEY),
       ('SecurityInformation', SECURITY_INFORMATION),
       ('pRpcSecurityDescriptor', RPC_SECURITY_DESCRIPTOR),
    }

 type BaseRegSetKeySecurityResponse struct { // NDRCALL: (
       ('ErrorCode', error_status_t),
    }

// 3.1.5.22 BaseRegSetValue (Opnum 22)
 type BaseRegSetValue struct { // NDRCALL:
    opnum = 22 (
       ('hKey', RPC_HKEY),
       ('lpValueName', RRP_UNICODE_STRING),
       ('dwType', DWORD),
       ('lpData', NDRUniConformantArray),
       ('cbData', DWORD),
    }

 type BaseRegSetValueResponse struct { // NDRCALL: (
       ('ErrorCode', error_status_t),
    }

// 3.1.5.23 BaseRegUnLoadKey (Opnum 23)
 type BaseRegUnLoadKey struct { // NDRCALL:
    opnum = 23 (
       ('hKey', RPC_HKEY),
       ('lpSubKey', RRP_UNICODE_STRING),
    }

 type BaseRegUnLoadKeyResponse struct { // NDRCALL: (
       ('ErrorCode', error_status_t),
    }

// 3.1.5.24 BaseRegGetVersion (Opnum 26)
 type BaseRegGetVersion struct { // NDRCALL:
    opnum = 26 (
       ('hKey', RPC_HKEY),
    }

 type BaseRegGetVersionResponse struct { // NDRCALL: (
       ('lpdwVersion', DWORD),
       ('ErrorCode', error_status_t),
    }

// 3.1.5.25 OpenCurrentConfig (Opnum 27)
 type OpenCurrentConfig struct { // NDRCALL:
    opnum = 27 (
       ('ServerName', PREGISTRY_SERVER_NAME),
       ('samDesired', REGSAM),
    }

 type OpenCurrentConfigResponse struct { // NDRCALL: (
       ('phKey', RPC_HKEY),
       ('ErrorCode', error_status_t),
    }

// 3.1.5.26 BaseRegQueryMultipleValues (Opnum 29)
 type BaseRegQueryMultipleValues struct { // NDRCALL:
    opnum = 29 (
       ('hKey', RPC_HKEY),
       ('val_listIn', RVALENT_ARRAY),
       ('num_vals', DWORD),
       ('lpvalueBuf', PBYTE_ARRAY),
       ('ldwTotsize', DWORD),
    }

 type BaseRegQueryMultipleValuesResponse struct { // NDRCALL: (
       ('val_listOut', RVALENT_ARRAY),
       ('lpvalueBuf', PBYTE_ARRAY),
       ('ldwTotsize', DWORD),
       ('ErrorCode', error_status_t),
    }

// 3.1.5.27 BaseRegSaveKeyEx (Opnum 31)
 type BaseRegSaveKeyEx struct { // NDRCALL:
    opnum = 31 (
       ('hKey', RPC_HKEY),
       ('lpFile', RRP_UNICODE_STRING),
       ('pSecurityAttributes', PRPC_SECURITY_ATTRIBUTES),
       ('Flags', DWORD),
    }

 type BaseRegSaveKeyExResponse struct { // NDRCALL: (
       ('ErrorCode', error_status_t),
    }

// 3.1.5.28 OpenPerformanceText (Opnum 32)
 type OpenPerformanceText struct { // NDRCALL:
    opnum = 32 (
       ('ServerName', PREGISTRY_SERVER_NAME),
       ('samDesired', REGSAM),
    }

 type OpenPerformanceTextResponse struct { // NDRCALL: (
       ('phKey', RPC_HKEY),
       ('ErrorCode', error_status_t),
    }

// 3.1.5.29 OpenPerformanceNlsText (Opnum 33)
 type OpenPerformanceNlsText struct { // NDRCALL:
    opnum = 33 (
       ('ServerName', PREGISTRY_SERVER_NAME),
       ('samDesired', REGSAM),
    }

 type OpenPerformanceNlsTextResponse struct { // NDRCALL: (
       ('phKey', RPC_HKEY),
       ('ErrorCode', error_status_t),
    }

// 3.1.5.30 BaseRegQueryMultipleValues2 (Opnum 34)
 type BaseRegQueryMultipleValues2 struct { // NDRCALL:
    opnum = 34 (
       ('hKey', RPC_HKEY),
       ('val_listIn', RVALENT_ARRAY),
       ('num_vals', DWORD),
       ('lpvalueBuf', PBYTE_ARRAY),
       ('ldwTotsize', DWORD),
    }

 type BaseRegQueryMultipleValues2Response struct { // NDRCALL: (
       ('val_listOut', RVALENT_ARRAY),
       ('lpvalueBuf', PBYTE_ARRAY),
       ('ldwRequiredSize', DWORD),
       ('ErrorCode', error_status_t),
    }

// 3.1.5.31 BaseRegDeleteKeyEx (Opnum 35)
 type BaseRegDeleteKeyEx struct { // NDRCALL:
    opnum = 35 (
       ('hKey', RPC_HKEY),
       ('lpSubKey', RRP_UNICODE_STRING),
       ('AccessMask', REGSAM),
       ('Reserved', DWORD),
    }

 type BaseRegDeleteKeyExResponse struct { // NDRCALL: (
       ('ErrorCode', error_status_t),
    }

//###############################################################################
// OPNUMs and their corresponding structures
//###############################################################################
OPNUMS = {
 0 : (OpenClassesRoot, OpenClassesRootResponse),
 1 : (OpenCurrentUser, OpenCurrentUserResponse),
 2 : (OpenLocalMachine, OpenLocalMachineResponse),
 3 : (OpenPerformanceData, OpenPerformanceDataResponse),
 4 : (OpenUsers, OpenUsersResponse),
 5 : (BaseRegCloseKey, BaseRegCloseKeyResponse),
 6 : (BaseRegCreateKey, BaseRegCreateKeyResponse),
 7 : (BaseRegDeleteKey, BaseRegDeleteKeyResponse),
 8 : (BaseRegDeleteValue, BaseRegDeleteValueResponse),
 9 : (BaseRegEnumKey, BaseRegEnumKeyResponse),
10 : (BaseRegEnumValue, BaseRegEnumValueResponse),
11 : (BaseRegFlushKey, BaseRegFlushKeyResponse),
12 : (BaseRegGetKeySecurity, BaseRegGetKeySecurityResponse),
13 : (BaseRegLoadKey, BaseRegLoadKeyResponse),
15 : (BaseRegOpenKey, BaseRegOpenKeyResponse),
16 : (BaseRegQueryInfoKey, BaseRegQueryInfoKeyResponse),
17 : (BaseRegQueryValue, BaseRegQueryValueResponse),
18 : (BaseRegReplaceKey, BaseRegReplaceKeyResponse),
19 : (BaseRegRestoreKey, BaseRegRestoreKeyResponse),
20 : (BaseRegSaveKey, BaseRegSaveKeyResponse),
21 : (BaseRegSetKeySecurity, BaseRegSetKeySecurityResponse),
22 : (BaseRegSetValue, BaseRegSetValueResponse),
23 : (BaseRegUnLoadKey, BaseRegUnLoadKeyResponse),
26 : (BaseRegGetVersion, BaseRegGetVersionResponse),
27 : (OpenCurrentConfig, OpenCurrentConfigResponse),
29 : (BaseRegQueryMultipleValues, BaseRegQueryMultipleValuesResponse),
31 : (BaseRegSaveKeyEx, BaseRegSaveKeyExResponse),
32 : (OpenPerformanceText, OpenPerformanceTextResponse),
33 : (OpenPerformanceNlsText, OpenPerformanceNlsTextResponse),
34 : (BaseRegQueryMultipleValues2, BaseRegQueryMultipleValues2Response),
35 : (BaseRegDeleteKeyEx, BaseRegDeleteKeyExResponse),
}

//###############################################################################
// HELPER FUNCTIONS
//###############################################################################
 func checkNullString(string interface{}){
    if string == NULL {
        return string

    if string[-1:] != '\x00' {
        return string + '\x00'
    } else  {
        return string

 func packValue(valueType, value interface{}){
    if valueType == REG_DWORD {
        retData = pack('<L', value)
    elif valueType == REG_DWORD_BIG_ENDIAN {
        retData = pack('>L', value)
    elif valueType == REG_EXPAND_SZ {
        try:
            retData = value.encode("utf-16le")
        except UnicodeDecodeError:
            import sys
            retData = value.decode(sys.getfilesystemencoding()).encode("utf-16le")
    elif valueType == REG_MULTI_SZ {
        try:
            retData = value.encode("utf-16le")
        except UnicodeDecodeError:
            import sys
            retData = value.decode(sys.getfilesystemencoding()).encode("utf-16le")
    elif valueType == REG_QWORD {
        retData = pack('<Q', value)
    elif valueType == REG_QWORD_LITTLE_ENDIAN {
        retData = pack('>Q', value)
    elif valueType == REG_SZ {
        try:
            retData = value.encode("utf-16le")
        except UnicodeDecodeError:
            import sys
            retData = value.decode(sys.getfilesystemencoding()).encode("utf-16le")
    } else  {
        retData = value

    return retData

 func unpackValue(valueType, value interface{}){
    if valueType == REG_DWORD {
        retData = unpack('<L', b''.join(value))[0]
    elif valueType == REG_DWORD_BIG_ENDIAN {
        retData = unpack('>L', ''.join(value))[0]
    elif valueType == REG_EXPAND_SZ {
        retData = "".join(value).decode("utf-16le")
    elif valueType == REG_MULTI_SZ {
        retData = "".join(value).decode("utf-16le")
    elif valueType == REG_QWORD {
        retData = unpack('<Q', ''.join(value))[0]
    elif valueType == REG_QWORD_LITTLE_ENDIAN {
        retData = unpack('>Q', ''.join(value))[0]
    elif valueType == REG_SZ {
        retData = b''.join(value).decode("utf-16le")
    } else  {
        retData = b''.join(value)

    return retData

 func hOpenClassesRoot(dce, samDesired = MAXIMUM_ALLOWED interface{}){
    request = OpenClassesRoot()
    request["ServerName"] = NULL
    request["samDesired"] = samDesired
    return dce.request(request)

 func hOpenCurrentUser(dce, samDesired = MAXIMUM_ALLOWED interface{}){
    request = OpenCurrentUser()
    request["ServerName"] = NULL
    request["samDesired"] = samDesired
    return dce.request(request)

 func hOpenLocalMachine(dce, samDesired = MAXIMUM_ALLOWED interface{}){
    request = OpenLocalMachine()
    request["ServerName"] = NULL
    request["samDesired"] = samDesired
    return dce.request(request)

 func hOpenPerformanceData(dce, samDesired = MAXIMUM_ALLOWED interface{}){
    request = OpenPerformanceData()
    request["ServerName"] = NULL
    request["samDesired"] = samDesired
    return dce.request(request)

 func hOpenUsers(dce, samDesired = MAXIMUM_ALLOWED interface{}){
    request = OpenUsers()
    request["ServerName"] = NULL
    request["samDesired"] = samDesired
    return dce.request(request)

 func hBaseRegCloseKey(dce, hKey interface{}){
    request = BaseRegCloseKey()
    request["hKey"] = hKey
    return dce.request(request)

 func hBaseRegCreateKey(dce, hKey, lpSubKey, lpClass = NULL, dwOptions = 0x00000001, samDesired = MAXIMUM_ALLOWED, lpSecurityAttributes = NULL, lpdwDisposition = REG_CREATED_NEW_KEY interface{}){
    request = BaseRegCreateKey()
    request["hKey"] = hKey
    request["lpSubKey"] = checkNullString(lpSubKey)
    request["lpClass"] = checkNullString(lpClass)
    request["dwOptions"] = dwOptions
    request["samDesired"] = samDesired
    if lpSecurityAttributes == NULL {
        request["lpSecurityAttributes"]["RpcSecurityDescriptor"]["lpSecurityDescriptor"] = NULL
    } else  {
        request["lpSecurityAttributes"] = lpSecurityAttributes
    request["lpdwDisposition"] = lpdwDisposition

    return dce.request(request)

 func hBaseRegDeleteKey(dce, hKey, lpSubKey interface{}){
    request = BaseRegDeleteKey()
    request["hKey"] = hKey
    request["lpSubKey"] = checkNullString(lpSubKey)
    return dce.request(request)

 func hBaseRegEnumKey(dce, hKey, dwIndex, lpftLastWriteTime = NULL interface{}){
    request = BaseRegEnumKey()
    request["hKey"] = hKey
    request["dwIndex"] = dwIndex
    request.fields["lpNameIn"].fields["MaximumLength"] = 1024
    request.fields["lpNameIn"].fields["Data"].fields["Data"].fields["MaximumCount"] = 1024//2
    request["lpClassIn"] = " "* 64
    request["lpftLastWriteTime"] = lpftLastWriteTime

    return dce.request(request)

 func hBaseRegEnumValue(dce, hKey, dwIndex, dataLen=256 interface{}){
    request = BaseRegEnumValue()
    request["hKey"] = hKey
    request["dwIndex"] = dwIndex
    retries = 1

    // We need to be aware the size might not be enough, so let's catch ERROR_MORE_DATA exception
    while true:
        try:
            // Only the maximum length field of the lpValueNameIn is used to determine the buffer length to be allocated
            // by the service. Specify a string with a zero length but maximum length set to the largest buffer size
            // needed to hold the value names.
            request.fields["lpValueNameIn"].fields["MaximumLength"] = dataLen*2
            request.fields["lpValueNameIn"].fields["Data"].fields["Data"].fields["MaximumCount"] = dataLen

            request["lpData"] = b' ' * dataLen
            request["lpcbData"] = dataLen
            request["lpcbLen"] = dataLen
            resp = dce.request(request)
        except DCERPCSessionError as e:
            if retries > 1 {
                LOG.debug("Too many retries when calling hBaseRegEnumValue, aborting")
                raise
            if e.get_error_code() == system_errors.ERROR_MORE_DATA {
                // We need to adjust the size
                retries +=1
                dataLen = e.get_packet()["lpcbData"]
                continue
            } else  {
                raise
        } else  {
            break

    return resp

 func hBaseRegFlushKey(dce, hKey interface{}){
    request = BaseRegFlushKey()
    request["hKey"] = hKey
    return dce.request(request)

 func hBaseRegGetKeySecurity(dce, hKey, securityInformation = OWNER_SECURITY_INFORMATION  interface{}){
    request = BaseRegGetKeySecurity()
    request["hKey"] = hKey
    request["SecurityInformation"] = securityInformation
    request["pRpcSecurityDescriptorIn"]["lpSecurityDescriptor"] = NULL
    request["pRpcSecurityDescriptorIn"]["cbInSecurityDescriptor"] = 1024

    return dce.request(request)

 func hBaseRegLoadKey(dce, hKey, lpSubKey, lpFile interface{}){
    request = BaseRegLoadKey()
    request["hKey"] = hKey
    request["lpSubKey"] = checkNullString(lpSubKey)
    request["lpFile"] = checkNullString(lpFile)
    return dce.request(request)

 func hBaseRegUnLoadKey(dce, hKey, lpSubKey interface{}){
    request = BaseRegUnLoadKey()
    request["hKey"] = hKey
    request["lpSubKey"] = checkNullString(lpSubKey)
    return dce.request(request)

 func hBaseRegOpenKey(dce, hKey, lpSubKey, dwOptions=0x00000001, samDesired = MAXIMUM_ALLOWED interface{}){
    request = BaseRegOpenKey()
    request["hKey"] = hKey
    request["lpSubKey"] = checkNullString(lpSubKey)
    request["dwOptions"] = dwOptions
    request["samDesired"] = samDesired 
    return dce.request(request)

 func hBaseRegQueryInfoKey(dce, hKey interface{}){
    request = BaseRegQueryInfoKey()
    request["hKey"] = hKey
    // Not the cleanest way, but oh well
    // Plus, Windows XP needs MaximumCount also set
    request.fields["lpClassIn"].fields["MaximumLength"] = 1024
    request.fields["lpClassIn"].fields["Data"].fields["Data"].fields["MaximumCount"] = 1024//2
    return dce.request(request)

 func hBaseRegQueryValue(dce, hKey, lpValueName, dataLen=512 interface{}){
    request = BaseRegQueryValue()
    request["hKey"] = hKey
    request["lpValueName"] = checkNullString(lpValueName)
    retries = 1

    // We need to be aware the size might not be enough, so let's catch ERROR_MORE_DATA exception
    while true:
        try:
            request["lpData"] =b' ' * dataLen
            request["lpcbData"] = dataLen
            request["lpcbLen"] = dataLen
            resp = dce.request(request)
        except DCERPCSessionError as e:
            if retries > 1 {
                LOG.debug("Too many retries when calling hBaseRegQueryValue, aborting")
                raise
            if e.get_error_code() == system_errors.ERROR_MORE_DATA {
                // We need to adjust the size
                dataLen = e.get_packet()["lpcbData"]
                continue
            } else  {
                raise
        } else  {
            break

    // Returns
    // ( dataType, data )
    return resp["lpType"], unpackValue(resp["lpType"], resp["lpData"])

 func hBaseRegReplaceKey(dce, hKey, lpSubKey, lpNewFile, lpOldFile interface{}){
    request = BaseRegReplaceKey()
    request["hKey"] = hKey
    request["lpSubKey"] = checkNullString(lpSubKey)
    request["lpNewFile"] = checkNullString(lpNewFile)
    request["lpOldFile"] = checkNullString(lpOldFile)
    return dce.request(request)

 func hBaseRegRestoreKey(dce, hKey, lpFile, flags=REG_REFRESH_HIVE interface{}){
    request = BaseRegRestoreKey()
    request["hKey"] = hKey
    request["lpFile"] = checkNullString(lpFile)
    request["Flags"] = flags
    return dce.request(request)

 func hBaseRegSaveKey(dce, hKey, lpFile, pSecurityAttributes = NULL interface{}){
    request = BaseRegSaveKey()
    request["hKey"] = hKey
    request["lpFile"] = checkNullString(lpFile)
    request["pSecurityAttributes"] = pSecurityAttributes
    return dce.request(request)

 func hBaseRegSetValue(dce, hKey, lpValueName, dwType, lpData interface{}){
    request = BaseRegSetValue()
    request["hKey"] = hKey
    request["lpValueName"] = checkNullString(lpValueName)
    request["dwType"] = dwType
    request["lpData"] = packValue(dwType,lpData)
    request["cbData"] = len(request["lpData"])
    return dce.request(request)

 func hBaseRegGetVersion(dce, hKey interface{}){
    request = BaseRegGetVersion()
    request["hKey"] = hKey
    return dce.request(request)

 func hOpenCurrentConfig(dce, samDesired = MAXIMUM_ALLOWED interface{}){
    request = OpenCurrentConfig()
    request["ServerName"] = NULL
    request["samDesired"] = samDesired
    return dce.request(request)

 func hBaseRegQueryMultipleValues(dce, hKey, val_listIn interface{}){
    // ToDo, check the result to see whether we need to 
    // have a bigger buffer for the data to receive
    request = BaseRegQueryMultipleValues()
    request["hKey"] = hKey

    for item in  val_listIn:
        itemn = RVALENT() 
        itemn["ve_valuename"] = checkNullString(item["ValueName"])
        itemn["ve_valuelen"] = len(itemn["ve_valuename"])
        itemn["ve_valueptr"] = NULL
        itemn["ve_type"] = item["ValueType"]
        request["val_listIn"].append(itemn)

    request["num_vals"] = len(request["val_listIn"])
    request["lpvalueBuf"] = list(b' '*128)
    request["ldwTotsize"] = 128

    resp = dce.request(request)
    retVal = list()
    for item in resp["val_listOut"]:
        itemn = dict()
        itemn["ValueName"] = item["ve_valuename"] 
        itemn["ValueData"] = unpackValue(item["ve_type"], resp["lpvalueBuf"][item["ve_valueptr"] : item["ve_valueptr"]+item["ve_valuelen"]])
        retVal.append(itemn)
 
    return retVal

 func hBaseRegSaveKeyEx(dce, hKey, lpFile, pSecurityAttributes = NULL, flags=1 interface{}){
    request = BaseRegSaveKeyEx()
    request["hKey"] = hKey
    request["lpFile"] = checkNullString(lpFile)
    request["pSecurityAttributes"] = pSecurityAttributes
    request["Flags"] = flags
    return dce.request(request)

 func hOpenPerformanceText(dce, samDesired = MAXIMUM_ALLOWED interface{}){
    request = OpenPerformanceText()
    request["ServerName"] = NULL
    request["samDesired"] = samDesired
    return dce.request(request)

 func hOpenPerformanceNlsText(dce, samDesired = MAXIMUM_ALLOWED interface{}){
    request = OpenPerformanceNlsText()
    request["ServerName"] = NULL
    request["samDesired"] = samDesired
    return dce.request(request)

 func hBaseRegDeleteValue(dce, hKey, lpValueName interface{}){
    request = BaseRegDeleteValue()
    request["hKey"] = hKey
    request["lpValueName"] = checkNullString(lpValueName)
    return dce.request(request)