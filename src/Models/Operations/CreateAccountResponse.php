<?php

/**
 * Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.
 */

declare(strict_types=1);

namespace formance\stack\Models\Operations;

use formance\stack\Models\Shared;
class CreateAccountResponse
{
    /**
     * HTTP response content type for this operation
     *
     * @var string $contentType
     */
    public string $contentType;

    /**
     * HTTP response status code for this operation
     *
     * @var int $statusCode
     */
    public int $statusCode;

    /**
     * Raw HTTP response; suitable for custom response parsing
     *
     * @var \Psr\Http\Message\ResponseInterface $rawResponse
     */
    public \Psr\Http\Message\ResponseInterface $rawResponse;

    /**
     * OK
     *
     * @var ?Shared\PaymentsAccountResponse $paymentsAccountResponse
     */
    public ?Shared\PaymentsAccountResponse $paymentsAccountResponse = null;

    /**
     * @param  ?string  $contentType
     * @param  ?int  $statusCode
     * @param  ?\Psr\Http\Message\ResponseInterface  $rawResponse
     * @param  ?Shared\PaymentsAccountResponse  $paymentsAccountResponse
     */
    public function __construct(?string $contentType = null, ?int $statusCode = null, ?\Psr\Http\Message\ResponseInterface $rawResponse = null, ?Shared\PaymentsAccountResponse $paymentsAccountResponse = null)
    {
        $this->contentType = $contentType;
        $this->statusCode = $statusCode;
        $this->rawResponse = $rawResponse;
        $this->paymentsAccountResponse = $paymentsAccountResponse;
    }
}