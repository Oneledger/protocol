provider "aws" {
  region     = "us-east-1"
}

resource "aws_vpc" "testNet" {
  cidr_block = "10.0.0.0/16"
  enable_dns_hostnames = true
  enable_dns_support = true
}

resource "aws_internet_gateway" "internetGW" {
  vpc_id = "${aws_vpc.testNet.id}"
}

resource "aws_route_table" "testNetRouteTable" {
  vpc_id = "${aws_vpc.testNet.id}"

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = "${aws_internet_gateway.internetGW.id}"
  }

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = "${aws_internet_gateway.internetGW.id}"
  }
}

resource "aws_network_acl" "main" {
  vpc_id = "${aws_vpc.testNet.id}"

  egress {
    protocol   = "all"
    rule_no    = 100
    action     = "allow"
    cidr_block = "0.0.0.0/0"
    from_port  = 0
    to_port    = 0
  }

  ingress {
    protocol   = "all"
    rule_no    = 100
    action     = "allow"
    cidr_block = "0.0.0.0/0"
    from_port  = 0
    to_port    = 0
  }
}

resource "aws_security_group" "security_groups_testNode" {
  name        = "ssh"
  description = "testNode security group"
  vpc_id      = "${aws_vpc.testNet.id}"

  ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_subnet" "subnet1" {
  vpc_id     = "${aws_vpc.testNet.id}"
  cidr_block = "10.0.16.0/20"
  availability_zone = "us-east-1a"
  map_public_ip_on_launch = true
}

resource "aws_subnet" "subnet2" {
  vpc_id     = "${aws_vpc.testNet.id}"
  cidr_block = "10.0.32.0/20"
  availability_zone = "us-east-1b"
  map_public_ip_on_launch = true
}

resource "aws_route_table_association" "a" {
  subnet_id      = "${aws_subnet.subnet1.id}"
  route_table_id = "${aws_route_table.testNetRouteTable.id}"
}

resource "aws_route_table_association" "b" {
  subnet_id      = "${aws_subnet.subnet2.id}"
  route_table_id = "${aws_route_table.testNetRouteTable.id}"
}

resource "aws_instance" "validatorNode" {
  count = 4
  instance_type = "t2.micro"
  ami = "ami-0ac019f4fcb7cb7e6"
  key_name = "stevie-one-ledger"
  availability_zone = "us-east-1a"
  subnet_id = "${aws_subnet.subnet1.id}"
  vpc_security_group_ids = ["${aws_security_group.security_groups_testNode.id}"]
  tags = {
    environment = "test"
  }
  
  provisioner "file" {
    source      = "bin"
    destination = "~/bin"

    connection {
      type     = "ssh"
      user     = "ubuntu"
      private_key = "${file("${aws_instance.validatorNode.0.key_name}.pem")}"
    }
  }

  provisioner "file" {
    source      = "config.toml"
    destination = "~/config.toml"

    connection {
      type     = "ssh"
      user     = "ubuntu"
      private_key = "${file("${aws_instance.validatorNode.0.key_name}.pem")}"
    }
  }
  
  provisioner "remote-exec" {
    script = "deploy.sh"
    connection {
      type     = "ssh"
      user     = "ubuntu"
      private_key = "${file("${aws_instance.validatorNode.0.key_name}.pem")}"
    }
  }
}

resource "aws_instance" "seedNode" {
  count = 2
  subnet_id = "${aws_subnet.subnet2.id}"
  instance_type = "t2.micro"
  ami = "ami-0ac019f4fcb7cb7e6"
  key_name = "stevie-one-ledger"
  availability_zone = "us-east-1b"
  vpc_security_group_ids = ["${aws_security_group.security_groups_testNode.id}"]
  tags = {
    environment = "test"
  }  
  provisioner "file" {
    source      = "bin"
    destination = "~/bin"

    connection {
      type     = "ssh"
      user     = "ubuntu"
      private_key = "${file("${aws_instance.seedNode.0.key_name}.pem")}"
    }
  }

  provisioner "file" {
    source      = "config.toml"
    destination = "~/config.toml"

    connection {
      type     = "ssh"
      user     = "ubuntu"
      private_key = "${file("${aws_instance.seedNode.0.key_name}.pem")}"
    }
  }
  
  provisioner "remote-exec" {
    script = "deploy.sh"
    connection {
      type     = "ssh"
      user     = "ubuntu"
      private_key = "${file("${aws_instance.seedNode.0.key_name}.pem")}"
    }
  }
}

resource "aws_elb" "seedNodeLB" {
  name               = "seed-node-lb-tf"
  internal           = false
  subnets = ["${aws_subnet.subnet1.id}","${aws_subnet.subnet2.id}"]
  security_groups = ["${aws_security_group.security_groups_testNode.id}"]
  instances                   = ["${aws_instance.seedNode.0.id}", "${aws_instance.seedNode.1.id}"]
  cross_zone_load_balancing   = false
  idle_timeout                = 400
  connection_draining         = true
  connection_draining_timeout = 400
  listener {
    instance_port     = 8000
    instance_protocol = "http"
    lb_port           = 80
    lb_protocol       = "http"
  }
  health_check {
    healthy_threshold   = 2
    unhealthy_threshold = 2
    timeout             = 3
    target              = "HTTP:8000/"
    interval            = 30
  }
  tags = {
    Environment = "production"
  }
}

output "public_ip" {
  value = "${join(",",aws_instance.validatorNode.*.public_ip)}"
}

output "public_dns" {
  value = "${join(",",aws_instance.validatorNode.*.public_dns)}"
}
