cmake_minimum_required(VERSION 2.7)
project(Paxos)

set(MINISTL_VERSION 0.1)
set(CMAKE_C_COMPILER g++)
set(EXECUTABLE_OUTPUT_PATH ../build/bin)
set(CMAKE_C_FLAGS ${CMAKE_C_FLAGS} -g -ggdb -Wall -std=c++11 )

add_executable(test ./Paxos.cpp ./src/Acceptor.cpp ./src/Proposer.cpp)
target_link_libraries(test pthread)