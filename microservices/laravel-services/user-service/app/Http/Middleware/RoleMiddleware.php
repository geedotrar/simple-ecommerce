<?php

namespace App\Http\Middleware;

use Closure;
use Illuminate\Http\Request;
use Tymon\JWTAuth\Facades\JWTAuth;

class RoleMiddleware
{
    public function handle(Request $request, Closure $next, ...$roles)
    {
        $user = JWTAuth::parseToken()->authenticate();
        $userRoles = $user->roles->pluck('name')->toArray();

        if (!array_intersect($roles, $userRoles)) {
            return response()->json(['error' => 'Forbidden'], 403);
        }

        return $next($request);
    }
}