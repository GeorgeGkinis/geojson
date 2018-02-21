
Vagrant.configure("2") do |config|

  config.vm.box = "centos/7"
  #config.vbguest.auto_update = false

  # Service-C
  config.vm.network "forwarded_port", guest: 8091, host: 8091, host_ip: "127.0.0.1"

  # minikube dashboard
  # config.vm.network "forwarded_port", guest: 30000, 30000

  config.vm.network "public_network"

    config.vm.provider "virtualbox" do |vb|
      vb.memory = 1024
      vb.cpus = 2
    end


  # Install Docker
  config.vm.provision "shell", path: "src/geojson/scripts/docker.install.sh"

  # Set file permissions
  config.vm.provision "shell", inline: "chmod 770 /vagrant/src/geojson/scripts/*"

  # Install Go
  config.vm.provision "shell", path: "src/geojson/scripts/go.install.sh"

  # Install Git
  config.vm.provision "shell", path: "src/geojson/scripts/git.install.sh"

  # Instal Kubectl
  # config.vm.provision "shell", path: "src/geojson/scripts/kubectl.install.sh"

  # Install Minikube
  # config.vm.provision "shell", path: "src/geojson/scripts/minikube.install.sh"

  # Enable password authentication
  config.vm.provision "shell", path: "src/geojson/scripts/enable.password.auth.sh"

  # Start containers
  config.vm.provision "shell", path: "src/geojson/scripts/docker.run.sh"

  # Reboot
  #config.vm.provision "shell", inline: "reboot"
end
