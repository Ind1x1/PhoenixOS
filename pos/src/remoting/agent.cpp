/*
 * Copyright 2025 The PhoenixOS Authors. All rights reserved.
 * 
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

#pragma once

#include <iostream>

#include "pos/include/command.h"
#include "pos/include/log.h"
#include "pos/include/agent.h"


POSAgent* pos_create_agent(){
    POSAgent *pos_agent = nullptr;
    POS_CHECK_POINTER(pos_agent = new POSAgent());
    return pos_agent;
}


int pos_destory_agent(POSAgent* pos_agent){
    POS_CHECK_POINTER(pos_agent);
    delete pos_agent;
    return 0;
}
