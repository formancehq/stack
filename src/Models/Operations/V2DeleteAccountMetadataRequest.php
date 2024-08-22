<?php

/**
 * Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.
 */

declare(strict_types=1);

namespace formance\stack\Models\Operations;

use formance\stack\Utils\SpeakeasyMetadata;
class V2DeleteAccountMetadataRequest
{
    /**
     * Account address
     *
     * @var string $address
     */
    #[SpeakeasyMetadata('pathParam:style=simple,explode=false,name=address')]
    public string $address;

    /**
     * The key to remove.
     *
     * @var string $key
     */
    #[SpeakeasyMetadata('pathParam:style=simple,explode=false,name=key')]
    public string $key;

    /**
     * Name of the ledger.
     *
     * @var string $ledger
     */
    #[SpeakeasyMetadata('pathParam:style=simple,explode=false,name=ledger')]
    public string $ledger;

    /**
     * @param  ?string  $address
     * @param  ?string  $key
     * @param  ?string  $ledger
     */
    public function __construct(?string $address = null, ?string $key = null, ?string $ledger = null)
    {
        $this->address = $address;
        $this->key = $key;
        $this->ledger = $ledger;
    }
}