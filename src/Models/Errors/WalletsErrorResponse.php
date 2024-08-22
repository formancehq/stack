<?php

/**
 * Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.
 */

declare(strict_types=1);

namespace formance\stack\Models\Errors;


/** WalletsErrorResponse - Error */
class WalletsErrorResponse
{
    /**
     *
     * @var SchemasWalletsErrorResponseErrorCode $errorCode
     */
    #[\JMS\Serializer\Annotation\SerializedName('errorCode')]
    #[\JMS\Serializer\Annotation\Type('\formance\stack\Models\Errors\SchemasWalletsErrorResponseErrorCode')]
    public SchemasWalletsErrorResponseErrorCode $errorCode;

    /**
     *
     * @var string $errorMessage
     */
    #[\JMS\Serializer\Annotation\SerializedName('errorMessage')]
    public string $errorMessage;

    /**
     * @param  ?SchemasWalletsErrorResponseErrorCode  $errorCode
     * @param  ?string  $errorMessage
     */
    public function __construct(?SchemasWalletsErrorResponseErrorCode $errorCode = null, ?string $errorMessage = null)
    {
        $this->errorCode = $errorCode;
        $this->errorMessage = $errorMessage;
    }
}