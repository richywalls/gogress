angular.module('goGress').controller('AppController', [
  '$scope',
  '$mdDialog',
  '$mdSidenav',
  '$mdToast',
  '$log',
  '$auth',
  function($scope, $mdDialog, $mdSidenav, $mdToast, $log, $auth) {
    $scope.sys_config = {};
    $scope.sys_config.font = "Coda"; //Roboto ???
    $scope.sys_config.font = "Roboto"; //Coda ???
    $scope.auth = $auth
    $scope.sys_config.allow_sidebar_right_locked_open = true;
    $scope.authenticate = function(provider) {
      $auth.authenticate(provider);
    }
    $scope.toggleLeft = function() {
      $mdSidenav('left').toggle()
        .then(function() {
          $log.debug("toggle left is done");
        });
    };
    $scope.toggleRight = function() {
      $mdSidenav('right').toggle()
        .then(function() {
          $log.debug("toggle right is done");
        });
    };
    $scope.closeRight = function() {
      if( $mdSidenav('right').isLockedOpen() ){
        $log.debug("toogle right is locked open");
        $scope.sys_config.allow_sidebar_right_locked_open = false;
      }
      $mdSidenav('right').close()
        .then(function() {
          $log.debug("toggle right has been closed");
        });
    };
    $scope.querySearch = false;
    $scope.enableSearch = function(call){
      $scope.querySearch = call;
    }
    $scope.disableSearch = function(){
      $scope.querySearch = false;
    }
    $scope.openToast = function(msg) {
      $mdToast.show(
        $mdToast.simple()
        .position("top right")
        .content(msg)
        .hideDelay(4000)
      );
    };

    $scope.showPictures = function($event) {
      $mdDialog.show({
        targetEvent: $event,
        templateUrl: 'partials/image-dialog.tpl.html',
        controller: 'ImageViewerController',
        controllerAs: 'imgViewerVm',
        locals: {
            portal: this.portal
        }
      });
    }

    $scope.copyToClipboard = function(text) {
      window.prompt("Copy to clipboard: Ctrl+C, Enter", text);
    }
  }
]);
