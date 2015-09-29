var services = angular.module('services', ['ngResource']);

services.factory('BOM', ['$resource', function($resource){
    return $resource('/api/bom', {}, {
      'reset': {method: 'POST', url: 'a/api/bom/reset'},
      'resetToSampleState': {method: 'POST', url: '/a/api/bom/resetToSampleBOM'}
    });
  }]);

services.factory('ORDERS', ['$resource', function($resource){
    return $resource('/api/orders/:id',{},{});
  }]);
