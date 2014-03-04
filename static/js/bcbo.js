angular.module('bcbo', ['ngRoute']).

factory('googJsApi', function($window, $q, $rootScope){
  var deferred = $q.defer();

  var doc = $window.document;
  var jsapi = doc.createElement('script');
  jsapi.src = 'https://www.google.com/jsapi';
  jsapi.onload = function() {
    $rootScope.$apply(function () {
      deferred.resolve($window.google);
    });
  };

  doc.body.appendChild(jsapi);
  return deferred.promise;
}).

factory('googChartApi', function($q, $rootScope, googJsApi){
  var deferred = $q.defer();
  var jsapi;

  googJsApi.then(function(google){
    jsapi = google;
    google.load('visualization', '1', {
      packages: ['corechart'],
      callback: onGoogChartApiLoaded
    });
  });

  function onGoogChartApiLoaded() {
    $rootScope.$apply(function(){
      deferred.resolve(jsapi.visualization);
    });
  };

  return deferred.promise;
}).

filter('truncate', function(){
  return function truncateFilter(value, max, tail) {
    if (!value) return;
    if (value.length <= max) return value;
    return value.substr(0, max) + (tail || '...');
  }
}).

directive('boChart', function(googChartApi){
  function boChartLinkFn(scope, elem, attrs) {
    scope.$watch('data', function(chartData){
      googChartApi.then(function(chartsApi){
        var data = chartData.rows;
        var header = [];
        for (var i=0, c; c = chartData.columns[i]; i++) {
          header.push(c.name);
        }
        data.unshift(header);
        data = chartsApi.arrayToDataTable(data);

        var opts = {
          title: chartData.title,
          isStacked: true,
          // chartArea: {height: '60%'},
          colors:['#F05940','#564F8A', '#00B5AD', '#6ECFF5', '#D95C5C', '#A1CF64', '#5C6166'],
          fontName: '"Open Sans", "Helvetica Neue", "Helvetica", "Arial", sans-serif',
          // theme: 'maximized'
          // legend: {
          //   position: 'in',
          //   alignment: 'end'
          // }
        };

        // var chart = new chartsApi.ColumnChart(elem[0]);
        var chart = new chartsApi[attrs.boChart](elem[0]);
        chart.draw(data, opts);
      });
    });
  };

  return {
    restrict: 'EA',
    scope: {bcChart: '=', data: '=chartData'},
    link: boChartLinkFn,
  }
}).

controller('UsersCtrl', function($scope, $http, $timeout){
  $scope.users = [];
  fetchUsersList();

  $scope.toggleExpand = function toggleExpand(user) {
    user.$expanded = !user.$expanded;
    if (user.$expanded) {
      fetchUserProfile(user);
      fetchUserCharts(user);
    }
  }

  function fetchUsersList() {
    $scope.users.$done = false;
    $http.get("/api/v1/xpeppers/users?offset=50&limit=10").
      then(function(resp){
        $scope.users.length = 0;
        $scope.users.push.apply($scope.users, resp.data.items);
        $scope.users.$done = true;
      });
  };

  function fetchUserProfile(user) {
    user.$fetched = false;
    user.profile = user.profile || {};
    $http.get("/api/v1/xpeppers/users/" + user.id + "/profile").
      then(function(resp){
        angular.copy(resp.data, user.profile);
        user.$fetched = true;
      });
  };

  function fetchUserCharts(user) {
    user.charts = user.charts || [];
    user.charts.$fetched = false;
    $http.get("/api/v1/xpeppers/users/" + user.id + "/charts").
      then(function(resp){
        user.charts.length = 0;
        user.charts.push.apply(user.charts, resp.data.charts);
        user.charts.$fetched = true;
      });
  }
}).

config(function($locationProvider, $routeProvider){
  $locationProvider.html5Mode(true);
  $routeProvider.
    when('/users', {
      templateUrl: '/static/users.tpl.html',
      controller: 'UsersCtrl',
      reloadOnSearch: false
    }).
    otherwise({
      redirectTo: '/users'
    });
});
