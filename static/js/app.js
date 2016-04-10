document.getElementById("demo").innerHTML = "My First JavaScript";

// var app = angular.module('myApp', [],
//     function($interpolateProvider){
//         $interpolateProvider.startSymbol('[[');
//         $interpolateProvider.endSymbol(']]');
//     });

var app = angular.module('myApp', []);

app.controller('myCtrl', function($scope) {
    $scope.firstName= "John";
    $scope.lastName= "Doe";
});
