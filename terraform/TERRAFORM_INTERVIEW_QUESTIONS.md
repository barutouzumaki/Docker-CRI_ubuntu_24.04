# Terraform Interview Questions

This document contains comprehensive interview questions for Terraform, covering basic concepts, advanced topics, real-world scenarios, and best practices.

---

## Table of Contents

1. [Basic Concepts](#basic-concepts)
2. [State Management](#state-management)
3. [Modules](#modules)
4. [Providers](#providers)
5. [Resources and Data Sources](#resources-and-data-sources)
6. [Variables and Outputs](#variables-and-outputs)
7. [Workspaces](#workspaces)
8. [Remote State](#remote-state)
9. [Best Practices](#best-practices)
10. [Troubleshooting](#troubleshooting)
11. [Advanced Topics](#advanced-topics)
12. [Real-World Scenarios](#real-world-scenarios)

---

## Basic Concepts

### 1. What is Terraform and why is it used?

**Answer:**

**Terraform** is an Infrastructure as Code (IaC) tool developed by HashiCorp that allows you to define, provision, and manage infrastructure using declarative configuration files.

**Key Features:**
- **Declarative Configuration**: Define desired state, not steps
- **Multi-Cloud Support**: Works with AWS, Azure, GCP, and many others
- **State Management**: Tracks infrastructure state
- **Plan and Apply**: Preview changes before applying
- **Idempotent**: Safe to run multiple times
- **Version Control**: Infrastructure code in Git

**Why Use Terraform:**
- **Consistency**: Same infrastructure across environments
- **Version Control**: Track infrastructure changes
- **Collaboration**: Team members can review and contribute
- **Reproducibility**: Recreate infrastructure easily
- **Documentation**: Code serves as documentation
- **Disaster Recovery**: Rebuild infrastructure from code

**Example:**
```hcl
provider "aws" {
  region = "us-east-1"
}

resource "aws_instance" "web" {
  ami           = "ami-0c55b159cbfafe1f0"
  instance_type = "t2.micro"
  
  tags = {
    Name = "WebServer"
  }
}
```

### 2. What is Infrastructure as Code (IaC)?

**Answer:**

**Infrastructure as Code (IaC)** is the practice of managing and provisioning infrastructure through machine-readable definition files, rather than through manual configuration or interactive configuration tools.

**Benefits:**
- **Version Control**: Track all infrastructure changes
- **Consistency**: Eliminate configuration drift
- **Speed**: Provision infrastructure quickly
- **Documentation**: Code documents infrastructure
- **Testing**: Test infrastructure changes
- **Disaster Recovery**: Rebuild from code
- **Cost Optimization**: Track and optimize resources

**Types of IaC:**
1. **Declarative** (Terraform, CloudFormation): Define desired state
2. **Imperative** (Ansible, Chef): Define steps to achieve state

### 3. What is the difference between Terraform and Ansible?

**Answer:**

| Aspect | Terraform | Ansible |
|--------|-----------|---------|
| **Primary Use** | Infrastructure provisioning | Configuration management |
| **Approach** | Declarative | Imperative/Declarative |
| **State Management** | Maintains state file | Stateless |
| **Idempotency** | Built-in | Requires idempotent modules |
| **Multi-Cloud** | Excellent | Good |
| **Agent** | Agentless | Agentless |
| **Best For** | Creating resources | Configuring existing resources |

**Terraform:**
- Creates and manages infrastructure
- Tracks state
- Best for provisioning

**Ansible:**
- Configures and manages systems
- No state tracking
- Best for configuration

**Often Used Together:**
- Terraform: Provision infrastructure
- Ansible: Configure provisioned infrastructure

### 4. What are the main components of Terraform?

**Answer:**

**1. Configuration Files (.tf):**
- Written in HCL (HashiCorp Configuration Language)
- Define resources, variables, outputs

**2. Providers:**
- Plugins that interact with APIs
- AWS, Azure, GCP, Kubernetes, etc.

**3. Resources:**
- Infrastructure components to create
- EC2 instances, S3 buckets, VPCs, etc.

**4. State File:**
- Tracks current infrastructure state
- Maps configuration to real resources

**5. Modules:**
- Reusable configuration blocks
- Encapsulate resources

**6. Variables:**
- Input parameters
- Make configurations flexible

**7. Outputs:**
- Exported values
- Share information between modules

**Example Structure:**
```
terraform/
├── main.tf          # Main configuration
├── variables.tf     # Variable definitions
├── outputs.tf       # Output definitions
├── terraform.tfstate # State file (local)
└── terraform.tfvars # Variable values
```

### 5. What is HCL (HashiCorp Configuration Language)?

**Answer:**

**HCL (HashiCorp Configuration Language)** is a configuration language created by HashiCorp for human-readable, machine-friendly configuration files.

**Features:**
- **Human-Readable**: Easy to write and understand
- **JSON-Compatible**: Can be converted to/from JSON
- **Syntax**: Similar to JSON but more flexible
- **Comments**: Supports single-line and multi-line comments
- **Interpolation**: Embed expressions in strings

**Example:**
```hcl
# Variable definition
variable "instance_type" {
  description = "EC2 instance type"
  type        = string
  default     = "t2.micro"
}

# Resource with interpolation
resource "aws_instance" "web" {
  ami           = "ami-0c55b159cbfafe1f0"
  instance_type = var.instance_type
  
  tags = {
    Name        = "WebServer-${var.environment}"
    Environment = var.environment
  }
}

# Output with function
output "instance_ip" {
  value = aws_instance.web.public_ip
  description = "Public IP of the instance"
}
```

**HCL vs JSON:**
- HCL: More readable, supports comments
- JSON: More strict, no comments
- Terraform can use both

### 6. What is the Terraform workflow?

**Answer:**

**Terraform Workflow:**

1. **Write Configuration:**
   ```hcl
   # main.tf
   resource "aws_instance" "web" {
     ami           = "ami-0c55b159cbfafe1f0"
     instance_type = "t2.micro"
   }
   ```

2. **Initialize:**
   ```bash
   terraform init
   ```
   - Downloads providers
   - Sets up backend
   - Initializes modules

3. **Format:**
   ```bash
   terraform fmt
   ```
   - Formats configuration files
   - Ensures consistent style

4. **Validate:**
   ```bash
   terraform validate
   ```
   - Validates syntax
   - Checks configuration

5. **Plan:**
   ```bash
   terraform plan
   ```
   - Shows execution plan
   - Preview changes
   - No modifications made

6. **Apply:**
   ```bash
   terraform apply
   ```
   - Creates/modifies resources
   - Updates state file

7. **Destroy (optional):**
   ```bash
   terraform destroy
   ```
   - Removes all resources
   - Cleans up infrastructure

**Best Practices:**
- Always run `terraform plan` before `apply`
- Review plan output
- Use version control
- Store state remotely
- Use workspaces for environments

### 7. What is terraform init?

**Answer:**

**`terraform init`** initializes a Terraform working directory by downloading providers and modules, and setting up the backend.

**What it does:**
1. **Downloads Providers:**
   - Downloads required provider plugins
   - Stores in `.terraform` directory

2. **Initializes Backend:**
   - Configures state storage
   - Sets up remote state if configured

3. **Downloads Modules:**
   - Downloads referenced modules
   - Local or remote modules

4. **Creates .terraform Directory:**
   - Stores provider binaries
   - Stores module cache

**Example:**
```bash
$ terraform init

Initializing the backend...

Initializing provider plugins...
- Finding latest version of hashicorp/aws...
- Installing hashicorp/aws v5.0.0...
- Installed hashicorp/aws v5.0.0

Terraform has been successfully initialized!
```

**When to run:**
- First time in a directory
- After adding new providers
- After changing backend configuration
- After adding new modules

**Flags:**
- `-upgrade`: Upgrade providers to latest versions
- `-reconfigure`: Reconfigure backend
- `-backend-config`: Backend configuration

### 8. What is terraform plan?

**Answer:**

**`terraform plan`** creates an execution plan showing what Terraform will do when you run `terraform apply`.

**What it does:**
- Reads current state
- Compares with configuration
- Shows what will be created, modified, or destroyed
- **Does NOT make any changes**

**Output:**
```
Terraform will perform the following actions:

  # aws_instance.web will be created
  + resource "aws_instance" "web" {
      + ami                          = "ami-0c55b159cbfafe1f0"
      + instance_type                = "t2.micro"
      + ...
    }

Plan: 1 to add, 0 to change, 0 to destroy.
```

**Use Cases:**
- Preview changes before applying
- Review infrastructure changes
- Share plans with team
- CI/CD pipelines

**Flags:**
- `-out=plan.out`: Save plan to file
- `-var`: Set variable values
- `-var-file`: Use variable file
- `-destroy`: Create destroy plan

**Example:**
```bash
# Create plan
terraform plan -out=tfplan

# Review plan
terraform show tfplan

# Apply saved plan
terraform apply tfplan
```

### 9. What is terraform apply?

**Answer:**

**`terraform apply`** executes the actions proposed in a Terraform plan to create, modify, or destroy infrastructure.

**What it does:**
1. Creates execution plan (if not provided)
2. Shows plan to user
3. Asks for confirmation
4. Applies changes
5. Updates state file

**Example:**
```bash
$ terraform apply

Terraform will perform the following actions:

  # aws_instance.web will be created
  + resource "aws_instance" "web" {
      + ami           = "ami-0c55b159cbfafe1f0"
      + instance_type = "t2.micro"
    }

Plan: 1 to add, 0 to change, 0 to destroy.

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

aws_instance.web: Creating...
aws_instance.web: Creation complete after 45s

Apply complete! Resources: 1 added, 0 changed, 0 destroyed.
```

**Flags:**
- `-auto-approve`: Skip confirmation prompt
- `-var`: Set variable values
- `-var-file`: Use variable file
- `-target`: Apply specific resource
- `-parallelism=n`: Limit concurrent operations

**Best Practices:**
- Always review plan first
- Use `-auto-approve` in CI/CD
- Use `-target` sparingly
- Monitor apply progress

### 10. What is terraform destroy?

**Answer:**

**`terraform destroy`** removes all infrastructure defined in the current configuration.

**What it does:**
- Identifies all resources in state
- Creates destroy plan
- Removes resources in reverse dependency order
- Updates state file

**Example:**
```bash
$ terraform destroy

Terraform will perform the following actions:

  # aws_instance.web will be destroyed
  - resource "aws_instance" "web" {
      - ami           = "ami-0c55b159cbfafe1f0" -> null
      - instance_type = "t2.micro" -> null
    }

Plan: 0 to add, 0 to change, 1 to destroy.

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

aws_instance.web: Destroying...
aws_instance.web: Destruction complete after 30s

Destroy complete! Resources: 1 destroyed.
```

**Use Cases:**
- Clean up test environments
- Remove temporary infrastructure
- Disaster recovery scenarios
- Cost optimization

**Flags:**
- `-auto-approve`: Skip confirmation
- `-target`: Destroy specific resource
- `-var`: Set variable values

**Best Practices:**
- Use with caution
- Always review destroy plan
- Backup state before destroy
- Use `-target` for selective destruction

### 11. What is the difference between terraform plan and terraform apply?

**Answer:**

| Aspect | terraform plan | terraform apply |
|--------|----------------|-----------------|
| **Purpose** | Preview changes | Execute changes |
| **Modifications** | None | Creates/modifies/destroys |
| **State** | Reads state | Reads and writes state |
| **Output** | Execution plan | Actual changes |
| **Use Case** | Review before apply | Make changes |

**terraform plan:**
- Shows what will happen
- No changes made
- Safe to run multiple times
- Use for review

**terraform apply:**
- Executes the plan
- Makes actual changes
- Updates state file
- Use to provision infrastructure

**Workflow:**
```bash
# 1. Plan first
terraform plan

# 2. Review plan output

# 3. Apply if satisfied
terraform apply

# Or save plan and apply later
terraform plan -out=tfplan
terraform apply tfplan
```

### 12. What is terraform fmt?

**Answer:**

**`terraform fmt`** formats Terraform configuration files to a canonical format and style.

**What it does:**
- Formats `.tf` and `.tfvars` files
- Ensures consistent formatting
- Fixes indentation and spacing
- Standardizes style

**Example:**
```bash
# Before formatting
resource "aws_instance" "web"{
ami="ami-0c55b159cbfafe1f0"
instance_type="t2.micro"
}

# After terraform fmt
resource "aws_instance" "web" {
  ami           = "ami-0c55b159cbfafe1f0"
  instance_type = "t2.micro"
}
```

**Usage:**
```bash
# Format all files in current directory
terraform fmt

# Format specific file
terraform fmt main.tf

# Check formatting without modifying
terraform fmt -check

# Recursive formatting
terraform fmt -recursive
```

**Best Practices:**
- Run before committing code
- Use in CI/CD pipelines
- Use `-check` flag in CI
- Format all files consistently

### 13. What is terraform validate?

**Answer:**

**`terraform validate`** validates the Terraform configuration files in a directory, checking for syntax errors and internal consistency.

**What it checks:**
- Syntax errors
- Type mismatches
- Missing required arguments
- Invalid references
- Provider configuration

**Example:**
```bash
$ terraform validate

Success! The configuration is valid.
```

**Error Example:**
```bash
$ terraform validate

Error: Missing required argument

  on main.tf line 5:
   5: resource "aws_instance" "web" {
  
The argument "ami" is required, but no definition was found.
```

**Usage:**
```bash
# Validate current directory
terraform validate

# Validate with variables
terraform validate -var="instance_type=t2.micro"
```

**Best Practices:**
- Run before plan/apply
- Use in CI/CD pipelines
- Fix validation errors early
- Validate after changes

### 14. What is terraform state?

**Answer:**

**Terraform state** is a file that tracks the mapping between your configuration and the real-world resources.

**Purpose:**
- Maps configuration to actual resources
- Tracks resource metadata
- Determines what needs to be created/modified/destroyed
- Stores resource dependencies

**State File Structure:**
```json
{
  "version": 4,
  "terraform_version": "1.5.0",
  "resources": [
    {
      "type": "aws_instance",
      "name": "web",
      "provider": "provider[\"registry.terraform.io/hashicorp/aws\"]",
      "instances": [
        {
          "attributes": {
            "ami": "ami-0c55b159cbfafe1f0",
            "id": "i-1234567890abcdef0",
            "instance_type": "t2.micro"
          }
        }
      ]
    }
  ]
}
```

**State Management:**
- Stored in `terraform.tfstate` (local) or remote backend
- Should be version controlled (with caution) or stored remotely
- Never edit manually
- Use `terraform state` commands to manage

**State Commands:**
```bash
# List resources in state
terraform state list

# Show resource details
terraform state show aws_instance.web

# Remove resource from state
terraform state rm aws_instance.web

# Move resource in state
terraform state mv aws_instance.old aws_instance.new
```

### 15. What is a Terraform provider?

**Answer:**

**A Terraform provider** is a plugin that Terraform uses to interact with APIs of cloud providers, SaaS providers, and other services.

**Common Providers:**
- **AWS**: `hashicorp/aws`
- **Azure**: `hashicorp/azurerm`
- **GCP**: `hashicorp/google`
- **Kubernetes**: `hashicorp/kubernetes`
- **Docker**: `kreuzwerker/docker`
- **GitHub**: `integrations/github`

**Provider Configuration:**
```hcl
# Configure AWS provider
provider "aws" {
  region  = "us-east-1"
  profile = "default"
  
  # Optional: version constraint
  version = "~> 5.0"
}

# Multiple provider instances
provider "aws" {
  alias  = "us-west"
  region = "us-west-2"
}

# Use alias
resource "aws_instance" "west" {
  provider = aws.us-west
  # ...
}
```

**Provider Requirements:**
```hcl
terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}
```

**Provider Features:**
- Resource management
- Data sources
- Authentication
- Region/endpoint configuration
- Version constraints

### 16. What is a Terraform resource?

**Answer:**

**A Terraform resource** is the most important element in Terraform configuration. It describes one or more infrastructure objects.

**Resource Syntax:**
```hcl
resource "resource_type" "resource_name" {
  # Resource arguments
  argument1 = value1
  argument2 = value2
  
  # Nested blocks
  block {
    nested_argument = value
  }
}
```

**Example:**
```hcl
resource "aws_instance" "web" {
  ami           = "ami-0c55b159cbfafe1f0"
  instance_type = "t2.micro"
  
  tags = {
    Name = "WebServer"
    Environment = "production"
  }
  
  # Lifecycle block
  lifecycle {
    create_before_destroy = true
  }
}
```

**Resource Types:**
- **Cloud Resources**: EC2, S3, VPC, etc.
- **Local Resources**: Files, directories
- **Null Resources**: For provisioners

**Resource Arguments:**
- Required arguments (must be provided)
- Optional arguments (have defaults)
- Computed arguments (set by provider)

**Resource Behavior:**
- **Create**: When resource doesn't exist
- **Update**: When configuration changes
- **Destroy**: When resource removed from config
- **Replace**: When certain changes require recreation

### 17. What is a Terraform data source?

**Answer:**

**A Terraform data source** allows you to fetch and use information from outside of Terraform or from another Terraform configuration.

**Purpose:**
- Query existing resources
- Get information not managed by Terraform
- Reference resources from other configurations
- Fetch external data

**Data Source Syntax:**
```hcl
data "data_source_type" "data_source_name" {
  # Query arguments
  filter {
    name   = "tag:Name"
    values = ["existing-instance"]
  }
}

# Use data source
resource "aws_instance" "new" {
  ami           = data.aws_ami.ubuntu.id
  instance_type = "t2.micro"
}
```

**Example:**
```hcl
# Get latest AMI
data "aws_ami" "ubuntu" {
  most_recent = true
  owners      = ["099720109477"]
  
  filter {
    name   = "name"
    values = ["ubuntu/images/hvm-ssd/ubuntu-jammy-22.04-amd64-server-*"]
  }
}

# Use AMI in resource
resource "aws_instance" "web" {
  ami           = data.aws_ami.ubuntu.id
  instance_type = "t2.micro"
}
```

**Common Data Sources:**
- `aws_ami`: Get AMI information
- `aws_vpc`: Get VPC information
- `aws_subnet`: Get subnet information
- `aws_security_group`: Get security group information
- `terraform_remote_state`: Get state from other configurations

**Data Source vs Resource:**
- **Data Source**: Read-only, fetches information
- **Resource**: Creates/manages infrastructure

### 18. What are Terraform variables?

**Answer:**

**Terraform variables** are input parameters that make your configuration flexible and reusable.

**Variable Declaration:**
```hcl
variable "instance_type" {
  description = "EC2 instance type"
  type        = string
  default     = "t2.micro"
  
  validation {
    condition     = contains(["t2.micro", "t2.small", "t2.medium"], var.instance_type)
    error_message = "Instance type must be t2.micro, t2.small, or t2.medium."
  }
}
```

**Variable Types:**
- `string`: Text value
- `number`: Numeric value
- `bool`: Boolean value
- `list`: List of values
- `map`: Key-value pairs
- `object`: Structured object
- `set`: Set of unique values
- `tuple`: Fixed-length list

**Variable Usage:**
```hcl
# In configuration
resource "aws_instance" "web" {
  ami           = "ami-0c55b159cbfafe1f0"
  instance_type = var.instance_type
}
```

**Setting Variables:**
```bash
# Command line
terraform apply -var="instance_type=t2.small"

# Variable file
terraform apply -var-file="production.tfvars"

# Environment variables
export TF_VAR_instance_type=t2.small
terraform apply
```

**Variable File (terraform.tfvars):**
```hcl
instance_type = "t2.small"
environment   = "production"
```

### 19. What are Terraform outputs?

**Answer:**

**Terraform outputs** expose values from your configuration for use by other Terraform configurations or external systems.

**Output Declaration:**
```hcl
output "instance_ip" {
  description = "Public IP of the instance"
  value       = aws_instance.web.public_ip
  sensitive   = false
}

output "instance_id" {
  value     = aws_instance.web.id
  sensitive = true
}
```

**Output Usage:**
```bash
# View all outputs
terraform output

# View specific output
terraform output instance_ip

# View as JSON
terraform output -json

# View sensitive outputs
terraform output -json | jq -r '.instance_id.value'
```

**Output in Other Configurations:**
```hcl
# Using remote state
data "terraform_remote_state" "network" {
  backend = "s3"
  config = {
    bucket = "my-terraform-state"
    key    = "network/terraform.tfstate"
  }
}

# Use output
resource "aws_instance" "web" {
  subnet_id = data.terraform_remote_state.network.outputs.public_subnet_id
}
```

**Output Features:**
- **Value**: The value to output
- **Description**: Documentation
- **Sensitive**: Hide value in console
- **Depends_on**: Explicit dependencies

**Best Practices:**
- Output important resource attributes
- Use descriptions
- Mark sensitive values
- Output IDs for references

### 20. What is the difference between variables and outputs?

**Answer:**

| Aspect | Variables | Outputs |
|--------|-----------|---------|
| **Direction** | Input to configuration | Output from configuration |
| **Purpose** | Make config flexible | Expose values |
| **Usage** | `var.variable_name` | `output.output_name` |
| **Set By** | User, files, environment | Computed from resources |
| **Scope** | Module input | Module output |

**Variables:**
- Input parameters
- Set by users
- Make configuration flexible
- Used with `var.variable_name`

**Outputs:**
- Exported values
- Computed from resources
- Share information
- Used with `output.output_name`

**Example:**
```hcl
# Input variable
variable "instance_type" {
  type    = string
  default = "t2.micro"
}

# Use variable
resource "aws_instance" "web" {
  instance_type = var.instance_type
  # ...
}

# Output value
output "instance_ip" {
  value = aws_instance.web.public_ip
}
```

---

*This is the first batch of 20 questions. More questions will be added in subsequent batches.*

