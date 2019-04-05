VAGRANTFILE_API_VERSION = '2'

Vagrant.configure(VAGRANTFILE_API_VERSION) do |config|
  config.vm.box = 'ubuntu/xenial32'
  config.vm.hostname = 'gophers'
  config.vm.network :forwarded_port, guest: 3000, host: 3000

  config.vm.provision :shell do |shell|
    shell.path = 'provision.sh'
    shell.privileged = false
  end

  config.ssh.forward_agent = true
  config.vm.box_check_update = true
end