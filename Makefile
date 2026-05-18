#
# Copyright (c) 2025-2026 TQ-Systems GmbH <license@tq-group.com>, D-82229 Seefeld, Germany. All rights reserved.
# Author: Christoph Krutz and the Energy Manager development team
#
# This software is licensed under the TQ-Systems Product Software License Agreement Version 1.0.3 or any later version.
# You can obtain a copy of the License Agreement in the TQS (TQ-Systems Software Licenses) folder on the following website:
# https://www.tq-group.com/en/support/downloads/tq-software-license-conditions/
# In case of any license issues please contact license@tq-group.com.
#

# The application ID must be unique across all Energy Manager applications. Ensure it does not
# conflict with other apps in the bundle. It is typically derived from the project URL. It must be
# hardcoded here because the official URL does not follow the common pattern for the app id.
# The app id in the bundle definition has to match the app id in the application.
APP_ID = go-demo

# Meta information about the application
APP_PRETTY_NAME = Go demo application
DESCRIPTION = This application is a demo application for Go.

# The application does not need valid timestamps
APPCLASS = no-time

include /opt/energy-manager/apps/Makefile

SERVICE_OPTS_SOCKET = --listenprotocol unix --listen ${SERVICE_RUN_DIR}/socket --listengroup www
