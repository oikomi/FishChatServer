var privateCloudStorageWebApp = angular.module("PrivateCloudStorageWeb", ['ngRoute',
	'LoginControllerModule', 'ControllersModule', 'ServicesModule']);

privateCloudStorageWebApp.run(function($rootScope) { 
	$rootScope.showLogin = true; 
	// console.log($rootScope)
}); 


privateCloudStorageWebApp.config(['$routeProvider', '$locationProvider', function($routeProvider, $locationProvider) { 
	//$locationProvider.html5Mode(true);
	$routeProvider 
		.when('/', { 
			templateUrl: 'views/login.html', 
			controller: 'LoginController' 
		}) 
		.when('/root', { 
			templateUrl: 'views/root.html', 
			controller: 'RootController' 
		}) 
		.otherwise({ 
			redirectTo: '/' 
		}); 
}]);