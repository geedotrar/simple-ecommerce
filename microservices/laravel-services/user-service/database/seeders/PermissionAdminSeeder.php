<?php

namespace Database\Seeders;

use App\Models\Permission;
use App\Models\Role;
use Illuminate\Database\Seeder;
use Illuminate\Support\Facades\Log;
use Illuminate\Support\Facades\Redis;

class PermissionAdminSeeder extends Seeder
{
    public function run(): void
    {
        Log::info('Starting admin permissions seeding');

        $adminRole = Role::firstOrCreate(['name' => 'admin']);

        $permissions = [
            'view_all_products',
            'view_active_products',
            'create_products',
            'update_products',
            'delete_products',
        ];

        $permissionIds = [];

        foreach ($permissions as $permName) {
            $permission = Permission::firstOrCreate(['name' => $permName]);
            $permissionIds[] = $permission->id;
        }

        $adminRole->permissions()->sync($permissionIds);
        
        Log::info('Admin permissions seeding completed');
    }
}
