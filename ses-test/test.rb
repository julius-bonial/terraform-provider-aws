#!/usr/bin/ruby

require 'rubygems'
require 'aws-sdk'
require 'yaml'

$region = 'us-east-1'


client = Aws::SES::Client.new(region: $region)

# The following example returns the notification attributes for an identity:

resp = client.get_identity_notification_attributes({
  identities: [
    "dev-stages.global", 
  ], 
})

puts resp.to_yaml
