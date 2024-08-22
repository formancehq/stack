<?php

/**
 * Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.
 */

declare(strict_types=1);

namespace formance\stack\Models\Operations;

use formance\stack\Utils\SpeakeasyMetadata;
class GetReconciliationRequest
{
    /**
     * The reconciliation ID.
     *
     * @var string $reconciliationID
     */
    #[SpeakeasyMetadata('pathParam:style=simple,explode=false,name=reconciliationID')]
    public string $reconciliationID;

    /**
     * @param  ?string  $reconciliationID
     */
    public function __construct(?string $reconciliationID = null)
    {
        $this->reconciliationID = $reconciliationID;
    }
}