#!/usr/bin/ruby
#
#
require 'json'

jsonText = gets.chomp

data = JSON.parse(jsonText)

if data.empty?()
  puts JSON.parse('{"name": "RubyTesting", "version": "v0.0.1", "requestVersion": ["v1"], "responseVersion":["v1"]}').to_json()
  return
end

d = data.dig("spec", "template", "spec", "nodeSelector", "beta.kubernetes.io/os")

response = {"version" => "v1", "isWhiteOut"=>false, "patches"=>[]}
if d.nil?
  puts response.to_json()
  return
end

# here we need to create a patch to remove beta.kubernetes.io and we need to add kubernetes.io and keep the value.
#

removalPatch = JSON.parse('{"op": "remove", "path": "/spec/template/spec/nodeSelector/beta.kubernetes.io~1os"}')
additionPatch = JSON.parse('{"op": "add", "path": "/spec/template/spec/nodeSelector/kubernetes.io~1os", "value":"'+d+'"}')

response['patches'] = response['patches'].append(removalPatch, additionPatch)

puts response.to_json
