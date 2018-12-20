require 'spec_helper'

RSpec.describe 'muchkeys command line', type: :aruba do
  before :each do
    Muchkeys.configure { |c|
      c.public_key = "#{__dir__}/../../../spec/data/tests.pem"
    }

    client = Muchkeys::ApplicationClient.new(Muchkeys.config)
    client.each_path do |path|
      client.allow_unsafe_operation do
        client.delete_key(path)
      end
    end
  end

  describe 'http' do
    it 'forwards http traffic' do
      run "punch http"
        stop_all_commands
      expect(last_command_started.stdout).to include 'localhost'
    end
  end

  describe 'https' do
    it 'forwards https traffic' do
      run "punch https"
      sleep 1
      stop_all_commands
      expect(last_command_started.stdout).to include 'localhost'
    end
  end

  describe 'version' do
    it 'spits out the version' do
      run 'punch --version'
      stop_all_commands
      expect(last_command_started.stdout).to include(Muchkeys::VERSION)
    end
  end
end