cmake_minimum_required(VERSION 2.8.9)
project(TMAC)
string(TOLOWER ${CMAKE_PROJECT_NAME} project_name_lowercase)

set(package_name "github.com/ztbrown/fdic")
set(devenv_image_name ${project_name_lowercase})
set(devenv_workpath "/gopath/src/${package_name}")
set(deployment_image_name ${project_name_lowercase})

set(deployment_image_version "1")

add_custom_command(
  OUTPUT  ${CMAKE_SOURCE_DIR}/vendor/vendor.stamp
  COMMAND docker run --rm --security-opt seccomp:unconfined  -t
  -v ${CMAKE_SOURCE_DIR}:${devenv_workpath}
  -w ${devenv_workpath}
  --name ${devenv_image_name}-devenv ${devenv_image_name}
  /bin/bash -c './scripts/vendor.sh'
  COMMAND cmake -E touch vendor/vendor.stamp
  WORKING_DIRECTORY ${CMAKE_SOURCE_DIR}
  DEPENDS go.mod devenv
  COMMENT "running go mod vendor"
  )

add_custom_target(
  vendor
  DEPENDS ${CMAKE_SOURCE_DIR}/vendor/vendor.stamp
  )

add_custom_target(devenv
  COMMAND docker build -t ${devenv_image_name}
  -f dockerfiles/devenv.Dockerfile ${CMAKE_SOURCE_DIR}
  WORKING_DIRECTORY ${CMAKE_SOURCE_DIR}
  COMMENT "building devenv docker image"
  )

add_custom_command(
  OUTPUT ${CMAKE_SOURCE_DIR}/bin/fdic
  COMMAND docker run --rm --security-opt seccomp:unconfined  -t
          -v ${CMAKE_SOURCE_DIR}:${devenv_workpath}
          -w ${devenv_workpath}
          --name ${devenv_image_name}-devenv ${devenv_image_name}
          /bin/bash -c './scripts/build.sh'
  WORKING_DIRECTORY ${CMAKE_SOURCE_DIR}
  DEPENDS devenv vendor 
  COMMENT "compiling project in devenv"
  )

add_custom_target(dockerize
  COMMAND docker build -t ${deployment_image_name}:${deployment_image_version}
  -f dockerfiles/deploy.Dockerfile ${CMAKE_SOURCE_DIR}
  WORKING_DIRECTORY ${CMAKE_SOURCE_DIR}
  DEPENDS ${CMAKE_SOURCE_DIR}/bin/fdic
  COMMENT "building deployment docker image"
  )


add_custom_target(fdic
  ALL
  DEPENDS ${CMAKE_SOURCE_DIR}/bin/fdic dockerize
  )
