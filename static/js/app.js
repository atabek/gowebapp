// var app = angular.module('myApp', [],
//     function($interpolateProvider){
//         $interpolateProvider.startSymbol('[[');
//         $interpolateProvider.endSymbol(']]');
//     });
var app = angular.module('aftercareApp', []);

app.controller('StudentCtrl', function($scope, $http) {
    'use strict';

    $scope.students = [];
    $scope.errors = [];

    $scope.search = function (row) {
        return (angular.lowercase(row.Student_id).indexOf(angular.lowercase($scope.query) || '') !== -1 ||
                angular.lowercase(row.Last_name).indexOf(angular.lowercase($scope.query)  || '') !== -1 ||
                angular.lowercase(row.First_name).indexOf(angular.lowercase($scope.query) || '') !== -1 ||
                angular.lowercase(row.Grade).indexOf(angular.lowercase($scope.query)      || '') !== -1);
    };

    $http.get('/students').then(function(res) {
        $scope.students = res.data;
        console.log($scope.students);
        $scope.refilter();
    }, function(msg) {
        $scope.log(msg.data);
    });

    $scope.log = function(msg) {
        $scope.errors.push(msg);
    };
});
