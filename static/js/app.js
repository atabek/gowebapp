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
    }, function(msg) {
        $scope.log(msg.data);
    });
}])

.controller('ClockinCtrl', ['$scope', '$http', 'ClockinService', '$location',
            function($scope, $http , ClockinService, $location){
    var url = $location.absUrl().split('/');
    var studentID = url[url.length - 1];
    $scope.studentID = studentID;
    $scope.Clockins = ClockinService.query({id: studentID});
    // console.log($scope.Clockins);
    // Clockins.$promise.then(function(data){
    //     console.log(data.length);
    //     console.log(data);
    //     data.splice(-2, 2);
    //     $scope.Clockins = data;
    //     console.log($scope.Clockins.length);
    // }, function(error){
    //     console.log('oopps ' + error);
    // });

    var tf = new Date();

    $scope.dateRange = {
        from: new Date(tf.setDate(1)),
        to:   new Date()
    };

    // Get the total time
    $scope.getTotal = function(){
        var total = 0;
        var clockins = $scope.Clockins;
        if(typeof clockins == "undefined"){
            return total;
        }
        for(var i = 0; i < clockins.length; i++){
            var clockin = clockins[i];
            total += clockin.TotalTime;
        }
        return total;
    }

    // Functions for clicking the row and action column
    $scope.clickRow = function(){
        console.log("ClickRow method");
    }

    $scope.clickAction = function(e){
        console.log("Click 2 method");
        e.stopPropagation();
    }
}])

.filter("dateRangeFilter", function(){
    return function(items, from, to, scope){
        var result = [];
        var sum    = 0;
        var day_count  = 0;
        if(typeof items == "undefined"){
            return result;
        }
        var df = parseDate(from);
        var dt = parseDate(to) + 86399999;
        for (var i = 0; i < items.length; i++){
            var tf = items[i].InAt * 1000,
                tt = items[i].OutAt * 1000;
            if (tf >= df && tt <= dt)  {
                result.push(items[i]);
                sum += items[i].TotalTime;
                day_count += 1;
            }
        }
        scope.sum = sum;
        scope.day_count = day_count;
        return result;
    };
})

.filter("formatDateFilter", function(){
    return function(item){
        if (typeof item == "number"){
            return formatDate(item);
        } else {
            return formatDate(0);
        }
    }
});

function parseDate(input) {
    if (typeof input == "undefined"){
        return new Date();
    } else {
        var day = input.getDate();
        var month = input.getMonth();
        var year = input.getFullYear();
        d = new Date(year, month, day);
        return d.getTime();
    }
}

function formatDate(input){
    var totalSec = input;
    var hours = parseInt( totalSec / 3600 ) % 24;
    var minutes = parseInt( totalSec / 60 ) % 60;
    var seconds = totalSec % 60;

    var result = (hours < 10 ? "0" + hours : hours)        + "h:" +
                 (minutes < 10 ? "0" + minutes : minutes)  + "m:" +
                 (seconds  < 10 ? "0" + seconds : seconds) + "s";
    return result;
}
