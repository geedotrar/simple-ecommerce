<?php

namespace Database\Seeders;

use App\Models\Permission;
use App\Models\Role;
use Illuminate\Database\Seeder;
use Illuminate\Support\Facades\Log;
use Illuminate\Support\Facades\Redis;

class PermissionUserSeeder extends Seeder
{
    public function run(): void
    {
        Log::info('Starting user permissions seeding');

        $userRole = Role::firstOrCreate(['name' => 'user']);

        $permissions = [
            'view_active_products',
        ];

        $permissionIds = [];

        foreach ($permissions as $permName) {
            $permission = Permission::firstOrCreate(['name' => $permName]);
            $permissionIds[] = $permission->id;
        }

        $userRole->permissions()->sync($permissionIds);

        Log::info('User permissions seeding completed');
    }
}
