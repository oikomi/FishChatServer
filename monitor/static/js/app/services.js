//
// Copyright 2014-2015 Hong Miao (miaohong@miaohong.org). All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

var servicesModule = angular.module("ServicesModule", []);

servicesModule.factory('loginService', function($http) {
	var url = 'api/v1/monitor'; 

	var runLoginRequest = function(reqParams, postData) { 
		return $http.post(url, postData, {params: reqParams, timeout:5000});
	}; 
	return { 
		events: function(reqParams, postData) { 
				return runLoginRequest(reqParams, postData);
			} 
	}; 
});

servicesModule.factory('rebootService', function($http) {
	var url = 'api/v1/monitor'; 

	var runRebootRequest = function(reqParams) { 
		return $http.get(url, {params: reqParams, timeout:5000});
	}; 
	return { 
		events: function(reqParams) { 
				return runRebootRequest(reqParams);
			} 
	}; 
});


servicesModule.factory('getServerDataService', function($http) {
	var url = 'api/v1/monitor'; 
	var runRequest = function(reqParams) { 
		return $http.get(url, {params: reqParams, timeout:5000});
	}; 
	return { 
		events: function(reqParams) { 
				return runRequest(reqParams);
			} 
	}; 
});