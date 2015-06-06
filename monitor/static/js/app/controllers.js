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

var loginControllerModule = angular.module("LoginControllerModule", ['ngCookies', 'ServicesModule']);

loginControllerModule.controller('LoginController', function($scope, $location, $cookies, loginService) { 
	//$scope.showLogin = true; 
	var reqParams = {
		action : "login"
	};
	
	$scope.login = function(user) {
		console.log(user.name)
		var postData = {
			username : user.name,
			password : user.password
		};
		loginService.events(reqParams, postData).success(function(data, status, headers, config) {
			console.log(data);
			if (data.status == "0") {
				$cookies.isLogin = true;
				$location.path("/root");
			} else {
				alert("登录失败！");
			}
		}).error(function(data, status, headers, config) {
			alert("登录失败！");
		});
	}; 
});


var controllersModule = angular.module("ControllersModule", ['ngCookies', 'ServicesModule']);

controllersModule.controller('RootController', function($scope, $location, $cookies, getServerDataService) { 
	var isLogin = $cookies.isLogin;
	
	if (!isLogin) {
		$location.path("/");
	}

	// var reqParams = {
		// action : "get_total_status"
	// };
	
	// getServerDataService.events(reqParams).success(function(data, status, headers, config) {
		// console.log(data);

	// }).error(function(data, status, headers, config) {
		// alert("获取服务器数据失败！");
	// });
});

controllersModule.controller('MsgServerController', function($scope, $location, $cookies, getServerDataService) { 
	
});

