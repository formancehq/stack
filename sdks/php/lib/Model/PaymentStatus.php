<?php
/**
 * PaymentStatus
 *
 * PHP version 7.4
 *
 * @category Class
 * @package  Formance
 * @author   OpenAPI Generator team
 * @link     https://openapi-generator.tech
 */

/**
 * Formance Stack API
 *
 * Open, modular foundation for unique payments flows  # Introduction This API is documented in **OpenAPI format**.  # Authentication Formance Stack offers one forms of authentication:   - OAuth2 OAuth2 - an open protocol to allow secure authorization in a simple and standard method from web, mobile and desktop applications. <SecurityDefinitions />
 *
 * The version of the OpenAPI document: develop
 * Contact: support@formance.com
 * Generated by: https://openapi-generator.tech
 * OpenAPI Generator version: 6.4.0
 */

/**
 * NOTE: This class is auto generated by OpenAPI Generator (https://openapi-generator.tech).
 * https://openapi-generator.tech
 * Do not edit the class manually.
 */

namespace Formance\Model;
use \Formance\ObjectSerializer;

/**
 * PaymentStatus Class Doc Comment
 *
 * @category Class
 * @package  Formance
 * @author   OpenAPI Generator team
 * @link     https://openapi-generator.tech
 */
class PaymentStatus
{
    /**
     * Possible values of this enum
     */
    public const PENDING = 'PENDING';

    public const ACTIVE = 'ACTIVE';

    public const TERMINATED = 'TERMINATED';

    public const FAILED = 'FAILED';

    public const SUCCEEDED = 'SUCCEEDED';

    public const CANCELLED = 'CANCELLED';

    /**
     * Gets allowable values of the enum
     * @return string[]
     */
    public static function getAllowableEnumValues()
    {
        return [
            self::PENDING,
            self::ACTIVE,
            self::TERMINATED,
            self::FAILED,
            self::SUCCEEDED,
            self::CANCELLED
        ];
    }
}


