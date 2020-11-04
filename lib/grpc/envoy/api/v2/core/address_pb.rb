# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: envoy/api/v2/core/address.proto

require 'google/protobuf'

require 'envoy/api/v2/core/socket_option_pb'
require 'google/protobuf/wrappers_pb'
require 'udpa/annotations/migrate_pb'
require 'udpa/annotations/status_pb'
Google::Protobuf::DescriptorPool.generated_pool.build do
  add_file("envoy/api/v2/core/address.proto", :syntax => :proto3) do
    add_message "envoy.api.v2.core.Pipe" do
      optional :path, :string, 1
      optional :mode, :uint32, 2
    end
    add_message "envoy.api.v2.core.SocketAddress" do
      optional :protocol, :enum, 1, "envoy.api.v2.core.SocketAddress.Protocol"
      optional :address, :string, 2
      optional :resolver_name, :string, 5
      optional :ipv4_compat, :bool, 6
      oneof :port_specifier do
        optional :port_value, :uint32, 3
        optional :named_port, :string, 4
      end
    end
    add_enum "envoy.api.v2.core.SocketAddress.Protocol" do
      value :TCP, 0
      value :UDP, 1
    end
    add_message "envoy.api.v2.core.TcpKeepalive" do
      optional :keepalive_probes, :message, 1, "google.protobuf.UInt32Value"
      optional :keepalive_time, :message, 2, "google.protobuf.UInt32Value"
      optional :keepalive_interval, :message, 3, "google.protobuf.UInt32Value"
    end
    add_message "envoy.api.v2.core.BindConfig" do
      optional :source_address, :message, 1, "envoy.api.v2.core.SocketAddress"
      optional :freebind, :message, 2, "google.protobuf.BoolValue"
      repeated :socket_options, :message, 3, "envoy.api.v2.core.SocketOption"
    end
    add_message "envoy.api.v2.core.Address" do
      oneof :address do
        optional :socket_address, :message, 1, "envoy.api.v2.core.SocketAddress"
        optional :pipe, :message, 2, "envoy.api.v2.core.Pipe"
      end
    end
    add_message "envoy.api.v2.core.CidrRange" do
      optional :address_prefix, :string, 1
      optional :prefix_len, :message, 2, "google.protobuf.UInt32Value"
    end
  end
end

module Envoy
  module Api
    module V2
      module Core
        Pipe = ::Google::Protobuf::DescriptorPool.generated_pool.lookup("envoy.api.v2.core.Pipe").msgclass
        SocketAddress = ::Google::Protobuf::DescriptorPool.generated_pool.lookup("envoy.api.v2.core.SocketAddress").msgclass
        SocketAddress::Protocol = ::Google::Protobuf::DescriptorPool.generated_pool.lookup("envoy.api.v2.core.SocketAddress.Protocol").enummodule
        TcpKeepalive = ::Google::Protobuf::DescriptorPool.generated_pool.lookup("envoy.api.v2.core.TcpKeepalive").msgclass
        BindConfig = ::Google::Protobuf::DescriptorPool.generated_pool.lookup("envoy.api.v2.core.BindConfig").msgclass
        Address = ::Google::Protobuf::DescriptorPool.generated_pool.lookup("envoy.api.v2.core.Address").msgclass
        CidrRange = ::Google::Protobuf::DescriptorPool.generated_pool.lookup("envoy.api.v2.core.CidrRange").msgclass
      end
    end
  end
end