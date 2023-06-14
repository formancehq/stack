# webhooks

### Available Operations

* [activateConfig](#activateconfig) - Activate one config
* [changeConfigSecret](#changeconfigsecret) - Change the signing secret of a config
* [deactivateConfig](#deactivateconfig) - Deactivate one config
* [deleteConfig](#deleteconfig) - Delete one config
* [getManyConfigs](#getmanyconfigs) - Get many configs
* [insertConfig](#insertconfig) - Insert a new config
* [testConfig](#testconfig) - Test one config

## activateConfig

Activate a webhooks config by ID, to start receiving webhooks to its endpoint.

### Example Usage

```php
<?php

declare(strict_types=1);
require_once 'vendor/autoload.php';

use \formance\stack\SDK;
use \formance\stack\Models\Shared\Security;
use \formance\stack\Models\Operations\ActivateConfigRequest;

$sdk = SDK::builder()
    ->build();

try {
    $request = new ActivateConfigRequest();
    $request->id = '4997257d-dfb6-445b-929c-cbe2ab182818';

    $response = $sdk->webhooks->activateConfig($request);

    if ($response->configResponse !== null) {
        // handle response
    }
} catch (Exception $e) {
    // handle exception
}
```

## changeConfigSecret

Change the signing secret of the endpoint of a webhooks config.

If not passed or empty, a secret is automatically generated.
The format is a random string of bytes of size 24, base64 encoded. (larger size after encoding)


### Example Usage

```php
<?php

declare(strict_types=1);
require_once 'vendor/autoload.php';

use \formance\stack\SDK;
use \formance\stack\Models\Shared\Security;
use \formance\stack\Models\Operations\ChangeConfigSecretRequest;
use \formance\stack\Models\Shared\ConfigChangeSecret;

$sdk = SDK::builder()
    ->build();

try {
    $request = new ChangeConfigSecretRequest();
    $request->configChangeSecret = new ConfigChangeSecret();
    $request->configChangeSecret->secret = 'V0bivxRWveaoz08afqjU6Ko/jwO0Cb+3';
    $request->id = '4997257d-dfb6-445b-929c-cbe2ab182818';

    $response = $sdk->webhooks->changeConfigSecret($request);

    if ($response->configResponse !== null) {
        // handle response
    }
} catch (Exception $e) {
    // handle exception
}
```

## deactivateConfig

Deactivate a webhooks config by ID, to stop receiving webhooks to its endpoint.

### Example Usage

```php
<?php

declare(strict_types=1);
require_once 'vendor/autoload.php';

use \formance\stack\SDK;
use \formance\stack\Models\Shared\Security;
use \formance\stack\Models\Operations\DeactivateConfigRequest;

$sdk = SDK::builder()
    ->build();

try {
    $request = new DeactivateConfigRequest();
    $request->id = '4997257d-dfb6-445b-929c-cbe2ab182818';

    $response = $sdk->webhooks->deactivateConfig($request);

    if ($response->configResponse !== null) {
        // handle response
    }
} catch (Exception $e) {
    // handle exception
}
```

## deleteConfig

Delete a webhooks config by ID.

### Example Usage

```php
<?php

declare(strict_types=1);
require_once 'vendor/autoload.php';

use \formance\stack\SDK;
use \formance\stack\Models\Shared\Security;
use \formance\stack\Models\Operations\DeleteConfigRequest;

$sdk = SDK::builder()
    ->build();

try {
    $request = new DeleteConfigRequest();
    $request->id = '4997257d-dfb6-445b-929c-cbe2ab182818';

    $response = $sdk->webhooks->deleteConfig($request);

    if ($response->statusCode === 200) {
        // handle response
    }
} catch (Exception $e) {
    // handle exception
}
```

## getManyConfigs

Sorted by updated date descending

### Example Usage

```php
<?php

declare(strict_types=1);
require_once 'vendor/autoload.php';

use \formance\stack\SDK;
use \formance\stack\Models\Shared\Security;
use \formance\stack\Models\Operations\GetManyConfigsRequest;

$sdk = SDK::builder()
    ->build();

try {
    $request = new GetManyConfigsRequest();
    $request->endpoint = 'https://example.com';
    $request->id = '4997257d-dfb6-445b-929c-cbe2ab182818';

    $response = $sdk->webhooks->getManyConfigs($request);

    if ($response->configsResponse !== null) {
        // handle response
    }
} catch (Exception $e) {
    // handle exception
}
```

## insertConfig

Insert a new webhooks config.

The endpoint should be a valid https URL and be unique.

The secret is the endpoint's verification secret.
If not passed or empty, a secret is automatically generated.
The format is a random string of bytes of size 24, base64 encoded. (larger size after encoding)

All eventTypes are converted to lower-case when inserted.


### Example Usage

```php
<?php

declare(strict_types=1);
require_once 'vendor/autoload.php';

use \formance\stack\SDK;
use \formance\stack\Models\Shared\Security;
use \formance\stack\Models\Shared\ConfigUser;

$sdk = SDK::builder()
    ->build();

try {
    $request = new ConfigUser();
    $request->endpoint = 'https://example.com';
    $request->eventTypes = [
        'TYPE1',
    ];
    $request->secret = 'V0bivxRWveaoz08afqjU6Ko/jwO0Cb+3';

    $response = $sdk->webhooks->insertConfig($request);

    if ($response->configResponse !== null) {
        // handle response
    }
} catch (Exception $e) {
    // handle exception
}
```

## testConfig

Test a config by sending a webhook to its endpoint.

### Example Usage

```php
<?php

declare(strict_types=1);
require_once 'vendor/autoload.php';

use \formance\stack\SDK;
use \formance\stack\Models\Shared\Security;
use \formance\stack\Models\Operations\TestConfigRequest;

$sdk = SDK::builder()
    ->build();

try {
    $request = new TestConfigRequest();
    $request->id = '4997257d-dfb6-445b-929c-cbe2ab182818';

    $response = $sdk->webhooks->testConfig($request);

    if ($response->attemptResponse !== null) {
        // handle response
    }
} catch (Exception $e) {
    // handle exception
}
```
