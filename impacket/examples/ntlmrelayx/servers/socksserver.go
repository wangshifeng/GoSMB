//!/usr/bin/env python
// SECUREAUTH LABS. Copyright 2018 SecureAuth Corporation. All rights reserved.
//
// This software is provided under under a slightly modified version
// of the Apache Software License. See the accompanying LICENSE file
// for more information.
//
// SOCKS proxy server/client
//
// Author:
//  Alberto Solino (@agsolino)
//
// Description:
//  A simple SOCKS server that proxy connection to relayed connections
//
// ToDo:
// [ ] Handle better the SOCKS specification (RFC1928), e.g. BIND
// [ ] Port handlers should be dynamically subscribed, and coded in another place. This will help coding
//     proxies for different protocols (e.g. MSSQL)
from __future__ import division
from __future__ import print_function
import socketserver
import socket
import time
import logging
from queue import Queue
from struct import unpack, pack
from threading import Timer, Thread

from impacket import LOG
from impacket.dcerpc.v5.enum import Enum
from impacket.structure import Structure

// Amount of seconds each socks plugin keep alive function will be called
// It is up to each plugin to send the keep alive to the target or not in every hit.
// In some cases (e.g. SMB) it is not needed to send a keep alive every 30 secs.
KEEP_ALIVE_TIMER = 30.0

 type enumItems struct { // Enum:
    NO_AUTHENTICATION = 0
    GSSAPI            = 1
    USER_PASS         = 2
    UNACCEPTABLE      = 0xFF

 type replyField struct { // Enum:
    SUCCEEDED             = 0
    SOCKS_FAILURE         = 1
    NOT_ALLOWED           = 2
    NETWORK_UNREACHABLE   = 3
    HOST_UNREACHABLE      = 4
    CONNECTION_REFUSED    = 5
    TTL_EXPIRED           = 6
    COMMAND_NOT_SUPPORTED = 7
    ADDRESS_NOT_SUPPORTED = 8

 type ATYP struct { // Enum:
    IPv4 = 1
    DOMAINNAME = 3
    IPv6 = 4

 type SOCKS5_GREETINGS struct { // Structure: (
        ('VER','B=5'),
        //('NMETHODS','B=0'),
        ('METHODS','B*B'),
    }


 type SOCKS5_GREETINGS_BACK struct { // Structure: (
        ('VER','B=5'),
        ('METHODS','B=0'),
    }

 type SOCKS5_REQUEST struct { // Structure: (
        ('VER','B=5'),
        ('CMD','B=0'),
        ('RSV','B=0'),
        ('ATYP','B=0'),
        ('PAYLOAD',':'),
    }

 type SOCKS5_REPLY struct { // Structure: (
        ('VER','B=5'),
        ('REP','B=5'),
        ('RSV','B=0'),
        ('ATYP','B=1'),
        ('PAYLOAD',':="AAAAA"'),
    }

 type SOCKS4_REQUEST struct { // Structure: (
        ('VER','B=4'),
        ('CMD','B=0'),
        ('PORT','>H=0'),
         ADDR [4]byte // ="
        ('PAYLOAD',':'),
    }

 type SOCKS4_REPLY struct { // Structure: (
        ('VER','B=0'),
        ('REP','B=0x5A'),
         RSV uint16 // =0
         RSV uint32 // =0
    }

activeConnections = Queue()

// Taken from https://stackoverflow.com/questions/474528/what-is-the-best-way-to-repeatedly-execute-a-function-every-x-seconds-in-python
// Thanks https://stackoverflow.com/users/624066/mestrelion
 type RepeatedTimer struct { // object:
   func (self TYPE) __init__(interval, function, *args, **kwargs interface{}){
    self._timer = nil
    self.interval = interval
    self.function = function
    self.args = args
    self.kwargs = kwargs
    self.is_running = false
    self.next_call = time.time()
    self.start()

   func (self TYPE) _run(){
    self.is_running = false
    self.start()
    self.function(*self.args, **self.kwargs)

   func (self TYPE) start(){
    if not self.is_running {
      self.next_call += self.interval
      self._timer = Timer(self.next_call - time.time(), self._run)
      self._timer.start()
      self.is_running = true

   func (self TYPE) stop(){
    self._timer.cancel()
    self.is_running = false

// Base  type for Relay Socks Servers for different protocols  struct { // SMB, MSSQL, etc
// Besides using this base  type you need to define one global variable when struct {
// writing a plugin for socksplugins:
// PLUGIN_CLASS = "<name of the  type for the plugin>" struct {
 type SocksRelay: struct {
    PLUGIN_NAME = "Base Plugin"
    // The plugin scheme, for automatic registration with relay servers
    // Should be specified in full caps, e.g. LDAP, HTTPS
    PLUGIN_SCHEME = ""

     func (self TYPE) __init__(targetHost, targetPort, socksSocket, activeRelays interface{}){
        self.targetHost = targetHost
        self.targetPort = targetPort
        self.socksSocket = socksSocket
        self.sessionData = activeRelays["data"]
        self.username = nil
        self.clientConnection = nil
        self.activeRelays = activeRelays

     func (self TYPE) initConnection(){
        // Here we do whatever is necessary to leave the relay ready for processing incoming connections
        raise RuntimeError("Virtual Function")

     func (self TYPE) skipAuthentication(){
        // Charged of bypassing any authentication attempt from the client
        raise RuntimeError("Virtual Function")

     func (self TYPE) tunnelConnection(){
        // Charged of tunneling the rest of the connection
        raise RuntimeError("Virtual Function")

    @staticmethod
     func (self TYPE) getProtocolPort(){
        // Should return the port this relay works against
        raise RuntimeError("Virtual Function")


 func keepAliveTimer(server interface{}){
    LOG.debug("KeepAlive Timer reached. Updating connections")

    for target in list(server.activeRelays.keys()):
        for port in list(server.activeRelays[target].keys()):
            // Now cycle through the users
            for user in list(server.activeRelays[target][port].keys()):
                if user != 'data' and user != 'scheme' {
                    // Let's call the keepAlive method for the handler to keep the connection alive
                    if server.activeRelays[target][port][user]["inUse"] is false {
                        LOG.debug('Calling keepAlive() for %s@%s:%s' % (user, target, port))
                        try:
                            server.activeRelays[target][port][user]["protocolClient"].keepAlive()
                        except Exception as e:
                            LOG.debug("Exception:",exc_info=true)
                            LOG.debug('SOCKS: %s' % str(e))
                            if str(e).find("Broken pipe") >= 0 or str(e).find("reset by peer") >=0 or \
                                            str(e).find("Invalid argument") >= 0 or str(e).find("Server not connected") >=0:
                                // Connection died, taking out of the active list
                                del (server.activeRelays[target][port][user])
                                if len(list(server.activeRelays[target][port].keys())) == 1 {
                                    del (server.activeRelays[target][port])
                                LOG.debug('Removing active relay for %s@%s:%s' % (user, target, port))
                    } else  {
                        LOG.debug('Skipping %s@%s:%s since it\'s being used at the moment' % (user, target, port))

 func activeConnectionsWatcher(server interface{}){
    while true:
        // This call blocks until there is data, so it doesn't loop endlessly
        target, port, scheme, userName, client, data = activeConnections.get()
        // ToDo: Careful. Dicts are not thread safe right?
        if (target in server.activeRelays) is not true {
            server.activeRelays[target] = {}
        if (port in server.activeRelays[target]) is not true {
            server.activeRelays[target][port] = {}

        if (userName in server.activeRelays[target][port]) is not true {
            LOG.info('SOCKS: Adding %s@%s(%s) to active SOCKS connection. Enjoy' % (userName, target, port))
            server.activeRelays[target][port][userName] = {}
            // This is the protocolClient. Needed because we need to access the killConnection from time to time.
            // Inside this instance, you have the session attribute pointing to the relayed session.
            server.activeRelays[target][port][userName]["protocolClient"] = client
            server.activeRelays[target][port][userName]["inUse"] = false
            server.activeRelays[target][port][userName]["data"] = data
            // Just for the CHALLENGE data, we're storing this general
            server.activeRelays[target][port]["data"] = data
            // Let's store the protocol scheme, needed be used later when trying to find the right socks relay server to use
            server.activeRelays[target][port]["scheme"] = scheme

            // Default values in case somebody asks while we're gettting the data
            server.activeRelays[target][port][userName]["isAdmin"] = "N/A"
            // Do we have admin access in this connection?
            try:
                LOG.debug("Checking admin status for user %s" % str(userName))
                isAdmin = client.isAdmin()
                server.activeRelays[target][port][userName]["isAdmin"] = isAdmin
            except Exception as e:
                // Method not implemented
                server.activeRelays[target][port][userName]["isAdmin"] = "N/A"
                pass
            LOG.debug("isAdmin returned: %s" % server.activeRelays[target][port][userName]["isAdmin"])
        } else  {
            LOG.info('Relay connection for %s at %s(%d) already exists. Discarding' % (userName, target, port))
            client.killConnection()

 func webService(server interface{}){
    from flask import Flask, jsonify

    app = Flask(__name__)

    log = logging.getLogger("werkzeug")
    log.setLevel(logging.ERROR)

    @app.route("/")
     func index(){
        print(server.activeRelays)
        return "Relays available: %s!" % (len(server.activeRelays))

    @app.route('/ntlmrelayx/api/v1.0/relays', methods=["GET"])
     func get_relays(){
        relays = []
        for target in server.activeRelays:
            for port in server.activeRelays[target]:
                for user in server.activeRelays[target][port]:
                    if user != 'data' and user != 'scheme' {
                        protocol = server.activeRelays[target][port]["scheme"]
                        isAdmin = server.activeRelays[target][port][user]["isAdmin"]
                        relays.append([protocol, target, user, isAdmin, str(port)])
        return jsonify(relays)

    @app.route('/ntlmrelayx/api/v1.0/relays', methods=["GET"])
     func get_info(relay interface{}){
        pass

    app.run(host='0.0.0.0', port=9090)

 type SocksRequestHandler struct { // socketserver.BaseRequestHandler:
     func (self TYPE) __init__(request, client_address, server interface{}){
        self.__socksServer = server
        self.__ip, self.__port = client_address
        self.__connSocket= request
        self.__socksVersion = 5
        self.targetHost = nil
        self.targetPort = nil
        self.__NBSession= nil
        socketserver.BaseRequestHandler.__init__(self, request, client_address, server)

     func (self TYPE) sendReplyError(error = replyField.CONNECTION_REFUSED interface{}){

        if self.__socksVersion == 5 {
            reply = SOCKS5_REPLY()
            reply["REP"] = error.value
        } else  {
            reply = SOCKS4_REPLY()
            if error.value != 0 {
                reply["REP"] = 0x5B
        return self.__connSocket.sendall(reply.getData())

     func (self TYPE) handle(){
        LOG.debug("SOCKS: New Connection from %s(%s)" % (self.__ip, self.__port))

        data = self.__connSocket.recv(8192)
        grettings = SOCKS5_GREETINGS_BACK(data)
        self.__socksVersion = grettings["VER"]

        if self.__socksVersion == 5 {
            // We need to answer back with a no authentication response. We're not dealing with auth for now
            self.__connSocket.sendall(str(SOCKS5_GREETINGS_BACK()))
            data = self.__connSocket.recv(8192)
            request = SOCKS5_REQUEST(data)
        } else  {
            // We're in version 4, we just received the request
            request = SOCKS4_REQUEST(data)

        // Let's process the request to extract the target to connect.
        // SOCKS5
        if self.__socksVersion == 5 {
            if request["ATYP"] == ATYP.IPv4.value {
                self.targetHost = socket.inet_ntoa(request["PAYLOAD"][:4])
                self.targetPort = unpack('>H',request["PAYLOAD"][4:])[0]
            elif request["ATYP"] == ATYP.DOMAINNAME.value {
                hostLength = unpack('!B',request["PAYLOAD"][0])[0]
                self.targetHost = request["PAYLOAD"][1:hostLength+1]
                self.targetPort = unpack('>H',request["PAYLOAD"][hostLength+1:])[0]
            } else  {
                LOG.error("No support for IPv6 yet!")
        // SOCKS4
        } else  {
            self.targetPort = request["PORT"]

            // SOCKS4a
            if request["ADDR"][:3] == "\x00\x00\x00" and request["ADDR"][3] != "\x00" {
                nullBytePos = request["PAYLOAD"].find("\x00")

                if nullBytePos == -1 {
                    LOG.error("Error while reading SOCKS4a header!")
                } else  {
                    self.targetHost = request["PAYLOAD"].split('\0', 1)[1][:-1]
            } else  {
                self.targetHost = socket.inet_ntoa(request["ADDR"])

        LOG.debug('SOCKS: Target is %s(%s)' % (self.targetHost, self.targetPort))

        if self.targetPort != 53 {
            // Do we have an active connection for the target host/port asked?
            // Still don't know the username, but it's a start
            if self.targetHost in self.__socksServer.activeRelays {
                if (self.targetPort in self.__socksServer.activeRelays[self.targetHost]) is not true {
                    LOG.error('SOCKS: Don\'t have a relay for %s(%s)' % (self.targetHost, self.targetPort))
                    self.sendReplyError(replyField.CONNECTION_REFUSED)
                    return
            } else  {
                LOG.error('SOCKS: Don\'t have a relay for %s(%s)' % (self.targetHost, self.targetPort))
                self.sendReplyError(replyField.CONNECTION_REFUSED)
                return

        // Now let's get into the loops
        if self.targetPort == 53 {
            // Somebody wanting a DNS request. Should we handle this?
            s = socket.socket()
            try:
                LOG.debug('SOCKS: Connecting to %s(%s)' %(self.targetHost, self.targetPort))
                s.connect((self.targetHost, self.targetPort))
            except Exception as e:
                LOG.debug("Exception:", exc_info=true)
                LOG.error('SOCKS: %s' %str(e))
                self.sendReplyError(replyField.CONNECTION_REFUSED)
                return

            if self.__socksVersion == 5 {
                reply = SOCKS5_REPLY()
                reply["REP"] = replyField.SUCCEEDED.value
                addr, port = s.getsockname()
                reply["PAYLOAD"] = socket.inet_aton(addr) + pack('>H', port)
            } else  {
                reply = SOCKS4_REPLY()

            self.__connSocket.sendall(reply.getData())

            while true:
                try:
                    data = self.__connSocket.recv(8192)
                    if data == b'' {
                        break
                    s.sendall(data)
                    data = s.recv(8192)
                    self.__connSocket.sendall(data)
                except Exception as e:
                    LOG.debug("Exception:", exc_info=true)
                    LOG.error('SOCKS: %s', str(e))

        // Let's look if there's a relayed connection for our host/port
        scheme = nil
        if self.targetHost in self.__socksServer.activeRelays {
            if self.targetPort in self.__socksServer.activeRelays[self.targetHost] {
                scheme = self.__socksServer.activeRelays[self.targetHost][self.targetPort]["scheme"]

        if scheme is not nil {
            LOG.debug('Handler for port %s found %s' % (self.targetPort, self.__socksServer.socksPlugins[scheme]))
            relay = self.__socksServer.socksPlugins[scheme](self.targetHost, self.targetPort, self.__connSocket,
                                  self.__socksServer.activeRelays[self.targetHost][self.targetPort])

            try:
                relay.initConnection()

                // Let's answer back saying we've got the connection. Data is fake
                if self.__socksVersion == 5 {
                    reply = SOCKS5_REPLY()
                    reply["REP"] = replyField.SUCCEEDED.value
                    addr, port = self.__connSocket.getsockname()
                    reply["PAYLOAD"] = socket.inet_aton(addr) + pack('>H', port)
                } else  {
                    reply = SOCKS4_REPLY()

                self.__connSocket.sendall(reply.getData())

                if relay.skipAuthentication() is not true {
                    // Something didn't go right
                    // Close the socket
                    self.__connSocket.close()
                    return

                // Ok, so we have a valid connection to play with. Let's lock it while we use it so the Timer doesn't send a
                // keep alive to this one.
                self.__socksServer.activeRelays[self.targetHost][self.targetPort][relay.username]["inUse"] = true

                relay.tunnelConnection()
            except Exception as e:
                LOG.debug("Exception:", exc_info=true)
                LOG.debug('SOCKS: %s' % str(e))
                if str(e).find("Broken pipe") >= 0 or str(e).find("reset by peer") >=0 or \
                                str(e).find("Invalid argument") >= 0:
                    // Connection died, taking out of the active list
                    del(self.__socksServer.activeRelays[self.targetHost][self.targetPort][relay.username])
                    if len(list(self.__socksServer.activeRelays[self.targetHost][self.targetPort].keys())) == 1 {
                        del(self.__socksServer.activeRelays[self.targetHost][self.targetPort])
                    LOG.debug('Removing active relay for %s@%s:%s' % (relay.username, self.targetHost, self.targetPort))
                    self.sendReplyError(replyField.CONNECTION_REFUSED)
                    return
                pass

            // Freeing up this connection
            if relay.username is not nil {
                self.__socksServer.activeRelays[self.targetHost][self.targetPort][relay.username]["inUse"] = false
        } else  {
            LOG.error('SOCKS: I don\'t have a handler for this port')

        LOG.debug("SOCKS: Shutting down connection")
        try:
            self.sendReplyError(replyField.CONNECTION_REFUSED)
        except Exception as e:
            LOG.debug('SOCKS END: %s' % str(e))


 type SOCKS struct { // socketserver.ThreadingMixIn, socketserver.TCPServer:
     func (self TYPE) __init__(server_address=('0.0.0.0', 1080), handler_class=SocksRequestHandler interface{}){
        LOG.info('SOCKS proxy started. Listening at port %d', server_address[1] )

        self.activeRelays = {}
        self.socksPlugins = {}
        self.restAPI = nil
        self.activeConnectionsWatcher = nil
        self.supportedSchemes = []
        socketserver.TCPServer.allow_reuse_address = true
        socketserver.TCPServer.__init__(self, server_address, handler_class)

        // Let's register the socksplugins plugins we have
        from impacket.examples.ntlmrelayx.servers.socksplugins import SOCKS_RELAYS

        for relay in SOCKS_RELAYS:
            LOG.info('%s loaded..' % relay.PLUGIN_NAME)
            self.socksPlugins[relay.PLUGIN_SCHEME] = relay
            self.supportedSchemes.append(relay.PLUGIN_SCHEME)

        // Let's create a timer to keep the connections up.
        self.__timer = RepeatedTimer(KEEP_ALIVE_TIMER, keepAliveTimer, self)

        // Let's start our RESTful API
        self.restAPI = Thread(target=webService, args=(self, ))
        self.restAPI.daemon = true
        self.restAPI.start()

        // Let's start out worker for active connections
        self.activeConnectionsWatcher = Thread(target=activeConnectionsWatcher, args=(self, ))
        self.activeConnectionsWatcher.daemon = true
        self.activeConnectionsWatcher.start()

     func (self TYPE) shutdown(){
        self.__timer.stop()
        del self.restAPI
        del self.activeConnectionsWatcher
        return socketserver.TCPServer.shutdown(self)

if __name__ == '__main__' {
    from impacket.examples import logger
    logger.init()
    s = SOCKS()
    s.serve_forever()
