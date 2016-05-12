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
            headers: {'Content-Type': 'application/json'}
        }
    });
})

.controller('StudentCtrl', ['$scope', '$http', function($scope, $http) {
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
        console.log($scope.students);
    }, function(msg) {
        $scope.log(msg.data);
    });
}])

.controller('ClockinCtrl', ['$scope', '$http', 'ClockinService', '$location',
            function($scope, $http , ClockinService, $location){
    var url = $location.absUrl().split('/');
    var studentID = url[url.length - 1]
    var Clockins = ClockinService.query({id: studentID});
    Clockins.$promise.then(function(data){
        console.log(typeof data);
        data = data.splice(-2, 2);
        $scope.Clockins = data;
    }, function(error){
        console.log('oopps ' + error);
    });

    var tf = new Date();

    $scope.dateRange = {
        from: new Date(tf.setDate(1)),
        to:   new Date()
    };
}])

.filter("dateRangeFilter", function(){
    return function(items, from, to){
        var result = [];
        if(typeof items == "undefined"){
            return result;
        }
        var df = parseDate(from);
        var dt = parseDate(to);
        console.log(df);
        // console.log(Object.keys(items).length);
        // console.log("is it an array: " + items[0] instanceof Array);
        for (var i = 0; i < items.length; i++){
            var tf = new Date(items[i].InAt * 1000),
                tt = new Date(items[i].OutAt * 1000);
            if (tf > df && tt < dt)  {
                result.push(items[i]);
            }
        }
        return result;
    };
});

function parseDate(input) {
    var day = input.getDate();
    var month = input.getMonth();
    var year = input.getFullYear();
    d = new Date(year, month, day);

    return d.getTime();
}
