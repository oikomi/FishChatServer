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