#!/usr/bin/env bash

# Exit on error
set -eo pipefail


################################
# Utilities section

NC='\033[0m'
RED='\033[0;31m'
YELLOW='\033[38;5;208m'
GREEN='\033[0;32m'

success(){
  printf "${GREEN}${1} ${NC}\n"
}

warn(){
  printf "${YELLOW}${1} ${NC}\n"
}

error(){
  printf "${RED}${1} ${NC}\n"
}

info(){
  printf "${NC}${1} ${NC}\n"
}

check_colima_present(){
  if \
    [[ -f "/usr/local/bin/colima" ]] || \
    [[ -f "/usr/local/sbin/colima" ]] || \
    [[ -f "/opt/homebrew/bin/colima" ]] || \
    [[ -f "/opt/homebrew/sbin/colima" ]] || \
    [[ $(brew ls --versions colima) != "" ]]
  then
    echo 1
  else
    echo 0
  fi
}

check_jq_present(){
  if [[ $(brew ls --versions jq) != "" ]]
  then
    echo 1
  else
    echo 0
  fi
}

# End utilities section

################################
# OS version verification
# Requires minimum OS version of 11.x
MIN_MAJOR_VERSION="11"
OS_VERSION=$(sw_vers -productVersion)
SEMVER_OS=(${OS_VERSION//./ })
MAJOR_OS_VERSION=${SEMVER_OS[0]}

if [[ $MAJOR_OS_VERSION -lt "$MIN_MAJOR_VERSION" ]]; then
  error "The minimum macOS version required is 11.x."
  error "You are currently on $OS_VERSION."
  error "Please upgrade your macOS before proceeding."
  info "For automatic OS upgrades, please refer to https://confluence.deliveryhero.com/display/DHITHelpdesk/macOS+-+Automatic+Upgrade+requirements"
  info "\n"
  exit
fi

# End OS version verification

################################
# Homebrew installation verification
if [[ "$(command -v brew)" == "" ]]; then
  error "Homebrew is necessary for installation to proceed."
  error "Please follow instructions at https://brew.sh/"
  exit
fi

# End Homebrew installation verification

################################
# Check Colima installation and provision if not found
COLIMA_PRESENT=$(check_colima_present)
if [[ "$COLIMA_PRESENT" == "1" ]]
then
  success "Colima is already present on your machine."
else
  # Install Colima using brew
  info "Colima will now be installed using Homebrew."
  warn "Note: This will also update Homebrew itself and any outdated packages."
  read -n1 -p "Proceed with installation (Y/n)? " answer
  case ${answer:0:1} in
    y|Y )
      # Proceed
      info "\\n"

      # Install Colima
      brew install colima

      # Check for successful installation
      if [[ $(brew ls --versions colima) != "" ]]
      then
        success "Colima was successfully installed\\n"
      else
        error "Something seems to have gone wrong. If you are unable to fix it on your own, please post a message on #docker-desktop-migration on slack for support."
        exit 1
      fi
    ;;
    * )
      # Continue with script without installing Colima
    ;;
  esac
fi

# End Colima installation

################################
# Ensure docker config file uses the right credential store and context
DOCKER_CONFIG_FILE="$HOME/.docker/config.json"
# Check if docker config file is present
if [[ -f "$DOCKER_CONFIG_FILE" ]]; then
  info "\n\nDocker configuration file was found at $DOCKER_CONFIG_FILE"
  info "Please ensure that it has the following lines in it"
  info "{"
  info "  \"credStore\": \"osxkeychain\","
  info "  \"currentContext\": \"colima\""
  info "}"

  read -n1 -p "Apply the required changes [Y/n]: " answer
  case ${answer:0:1} in
    y|Y )
      # Install `jq` via brew before proceeding
      JQ_PRESENT=$(check_jq_present)
      if [[ "$JQ_PRESENT" == "0" ]]; then
        read -p "This will install `jq`. Proceed [Y/n]? " answer
        case ${answer:0:1} in 
          y|Y )
            brew install jq  
          ;;
        esac
      fi

      info "\\nApplying changes"
      contents="$(/usr/bin/env jq '.credsStore = "osxkeychain" | .currentContext = "colima"' $DOCKER_CONFIG_FILE)"
      echo -E "${contents}" > $DOCKER_CONFIG_FILE

      success "Docker configuration file was updated\n"
    ;;
  esac
fi


# End docker config fixes

################################
# Add fix for "Cannot connect to the Docker daemon at unix:///var/run/docker.sock. Is the docker daemon running?"
# https://github.com/abiosoft/colima/blob/main/docs/FAQ.md#cannot-connect-to-the-docker-daemon-at-unixvarrundockersock-is-the-docker-daemon-running
# Also see https://github.com/GoogleContainerTools/skaffold/issues/7078#issuecomment-1180979589
# And https://deliveryhero.slack.com/archives/C04BHRT9T1P/p1673879019100899?thread_ts=1672839710.726989&cid=C04BHRT9T1P
DOCKER_HOST="unix:///${HOME}/.colima/default/docker.sock"

info "\n\nTo ensure compatibility for applications that do not recognize custom Docker contexts, "
info "ensure environment variable DOCKER_HOST points to Colima's Docker socket."
info "By default, this is located as $DOCKER_HOST."
info "Add this to your default shell's configuration.\n"

# Check for ZSHRC
if [[ -f "$HOME/.zshrc" ]]; then
  info "Found zsh configuration at $HOME/.zshrc"
  
  # Check to see if the variable exists
  case `grep -Fx "export DOCKER_HOST=\"unix://${HOME}/.colima/default/docker.sock\"" "$HOME/.zshrc" > /dev/null; echo $?` in
    0)
      # Variable is already present, do nothing
      info "Good job, this is already present in your zsh configuration."
    ;;
    1)
      # Variable is absent
      read -n1 -p "Apply the required changes for zsh [Y/n]: " answer
      case ${answer:0:1} in
        y|Y )
          echo -e "\n# Colima: Set DOCKER_HOST variable" >> "$HOME/.zshrc"
          echo -e "# See: https://github.com/abiosoft/colima/blob/main/docs/FAQ.md#cannot-connect-to-the-docker-daemon-at-unixvarrundockersock-is-the-docker-daemon-running" >> "$HOME/.zshrc"
          echo -e "export DOCKER_HOST=\"unix://${HOME}/.colima/default/docker.sock\"" >> "$HOME/.zshrc"
          success "\nAdded DOCKER_HOST variable to $HOME/.zshrc"
        ;;
      esac
    ;;
  esac

fi

# End fix

info "\nIn case you face further issues, please refer to the FAQ at https://docs.google.com/document/d/1kLDBSl2cvkgTRZwkHXtfvArSyujbXMX7tVS19NU6FqI/edit#heading=h.qtgfkz554pgt"
info "\nIf you need further support, please reach out to us on #docker-desktop-migration on Slack."
info "\n"

exit 0

