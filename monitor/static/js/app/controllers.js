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

