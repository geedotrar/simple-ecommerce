<?php

use App\Http\Controllers\AccessController;
use App\Http\Controllers\AuthController;
use Illuminate\Support\Facades\Route;

Route::post('/register', [AuthController::class, 'register']);
Route::post('/login', [AuthController::class, 'login']);
Route::middleware('auth:api')->get('/access', [AccessController::class, 'getPermissions']);
Route::middleware('auth:api')->get('/validate-token', function () {
    return response()->json(['valid' => true]);
});
