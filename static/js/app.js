angular.module('aftercareApp', ['ngResource' /*, 'ngRoute' */])

.config(function($interpolateProvider) {
  $interpolateProvider.startSymbol('[[');
  $interpolateProvider.endSymbol(']]');
})

angular.module('aftercareApp')
.factory('ClockinService', function($resource){
    console.log("Inside ClockinService");
    return $resource('http://localhost:3000/clockins/student/json/:id',
        {id: '@id'},
        {query: {
            isArray: true,
            method: 'GET',
            headers: {'X-Api-Secret': 'xxx', 'Authorization': 'xxx', 'Content-Type': 'application/x-www-form-urlencoded'}
        }
    });
})

.service('SharedProperties', function () {
    var studentID = 1;

    return {
        getStudentID: function () {
            return studentID;
        },
        setStudentID: function(value) {
            studentID = value;
        }
    };
})

.controller('StudentCtrl', function($scope, $http, ClockinService, SharedProperties) {
    'use strict';

    $scope.students = [];
    $scope.errors = [];
    $scope.clockins = [];

    $scope.search = function (row) {
        return (angular.lowercase(row.Student_id).indexOf(angular.lowercase($scope.query) || '') !== -1 ||
                angular.lowercase(row.Last_name).indexOf(angular.lowercase($scope.query)  || '') !== -1 ||
                angular.lowercase(row.First_name).indexOf(angular.lowercase($scope.query) || '') !== -1 ||
                angular.lowercase(row.Grade).indexOf(angular.lowercase($scope.query)      || '') !== -1);
    };

    $http.get('/students').then(function(res) {
        $scope.students = res.data;
        console.log(typeof $scope.students);
        console.log($scope.students);
    }, function(msg) {
        $scope.log(msg.data);
    });

    $scope.setDataForStudent = function(studentID){
        SharedProperties.setStudentID(studentID);
        console.log(SharedProperties.getStudentID());
    };
})

.controller('ClockinCtrl', function($scope, $http , ClockinService, SharedProperties,  $location){
    var url = $location.absUrl().split('/');
    var studentID = url[url.length - 1]
    console.log(studentID);
    var Clockins = ClockinService.query({id: studentID});
    Clockins.$promise.then(function(data){
        //$scope.Clockins = angular.toJson(data);
        data = data.splice(-2, 2);
        $scope.Clockins = data;
        console.log($scope.Clockins);
    });
});
