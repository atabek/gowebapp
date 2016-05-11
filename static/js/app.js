angular.module('aftercareApp', ['ngResource'])

.config(function($interpolateProvider) {
  $interpolateProvider.startSymbol('[[');
  $interpolateProvider.endSymbol(']]');
})

angular.module('aftercareApp').factory('ClockinService', function($resource){
    return $resource('http://localhost:3000/clockins/student/json/:id',
        {id: '@id'},
        {query: {
            isArray: true,
            method: 'GET',
            headers: {'X-Api-Secret': 'xxx', 'Authorization': 'xxx', 'Content-Type': 'application/x-www-form-urlencoded'}
        }
    });
})

.filter("dateRangeFilter", function(){

})

.controller('StudentCtrl', function($scope, $http) {
    'use strict';

    $scope.students = [];

    $scope.search = function (row) {
        return (angular.lowercase(row.Student_id).indexOf(angular.lowercase($scope.query) || '') !== -1 ||
                angular.lowercase(row.Last_name).indexOf(angular.lowercase($scope.query)  || '') !== -1 ||
                angular.lowercase(row.First_name).indexOf(angular.lowercase($scope.query) || '') !== -1 ||
                angular.lowercase(row.Grade).indexOf(angular.lowercase($scope.query)      || '') !== -1);
    };

    $http.get('/students').then(function(res) {
        $scope.students = res.data;
    }, function(msg) {
        $scope.log(msg.data);
    });
})

.controller('ClockinCtrl', function($scope, $http , ClockinService, $location){
    var url = $location.absUrl().split('/');
    var studentID = url[url.length - 1]
    var Clockins = ClockinService.query({id: studentID});
    Clockins.$promise.then(function(data){
        //$scope.Clockins = angular.toJson(data);
        data = data.splice(-2, 2);
        $scope.Clockins = data;
    });

    var tf = new Date();

    $scope.dateRange = {
        from: new Date(tf.setDate(tf.getDate() - tf.getDate() + 1)),
        to:   new Date()
    };
    console.log(typeof $scope.dateRange.from);
});
