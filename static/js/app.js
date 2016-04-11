document.getElementById("demo").innerHTML = "My First JavaScript";

// var app = angular.module('myApp', [],
//     function($interpolateProvider){
//         $interpolateProvider.startSymbol('[[');
//         $interpolateProvider.endSymbol(']]');
//     });

var app = angular.module('aftercareApp', []);

app.controller('StudentCtrl', function($scope, $http, $filter) {
    'use strict';

    $scope.firstName = "Atabek";
    $scope.lastName  = "Akbalaev";

    $scope.students = [];
    $scope.errors = [];
    $scope.search = {};
    $scope.filteredFields = ['Student_id', 'Last_name', 'First_name', 'Grade'];
    $scope.filteredStudents = [];

    $scope.refilter = function() {
        for (var f in $scope.filteredFields) {
            var field = $scope.filteredFields[f];
            if ($scope.search[field] === '') delete $scope.search[field];
        }
        $scope.filteredStudents = $filter('filter')($scope.students, $scope.search);
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
